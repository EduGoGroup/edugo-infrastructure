package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// scriptRunner ejecuta scripts de migración JavaScript usando goja
// y mapeando operaciones a llamadas del driver oficial de MongoDB.
type scriptRunner struct {
	ctx context.Context
	db  *mongo.Database
	vm  *goja.Runtime
}

func newScriptRunner(ctx context.Context, db *mongo.Database) *scriptRunner {
	return &scriptRunner{
		ctx: ctx,
		db:  db,
		vm:  goja.New(),
	}
}

func (r *scriptRunner) Run(script string) error {
	r.injectGlobals()

	if _, err := r.vm.RunString(script); err != nil {
		return err
	}

	return nil
}

func (r *scriptRunner) injectGlobals() {
	dbObj := r.newDBObject()
	_ = r.vm.Set("db", dbObj)

	// Alias para compatibilidad
	_ = r.vm.Set("ISODate", r.isoDate)
	_ = r.vm.Set("ObjectId", r.objectID)
	_ = r.vm.Set("print", r.print)

	console := r.vm.NewObject()
	_ = console.Set("log", r.print)
	_ = console.Set("error", r.print)
	_ = r.vm.Set("console", console)
}

func (r *scriptRunner) newDBObject() *goja.Object {
	props := map[string]goja.Value{
		"createCollection": r.vm.ToValue(r.createCollection),
		"dropCollection":   r.vm.ToValue(r.dropCollection),
		"collection":       r.vm.ToValue(r.collection),
		"getCollection":    r.vm.ToValue(r.collection),
		"runCommand":       r.vm.ToValue(r.runCommand),
	}

	dyn := &mongoDBObject{
		runner: r,
		props:  props,
	}

	return r.vm.NewDynamicObject(dyn)
}

func (r *scriptRunner) createCollection(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("createCollection requiere nombre de colección")))
	}

	name := call.Argument(0).String()
	command := bson.D{{Key: "create", Value: name}}

	if len(call.Arguments) > 1 {
		opts, err := r.valueToMap(call.Argument(1))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}
		command = append(command, mapToBsonD(opts)...)
	}

	if err := r.db.RunCommand(r.ctx, command).Err(); err != nil {
		panic(r.vm.NewGoError(err))
	}

	return goja.Undefined()
}

func (r *scriptRunner) dropCollection(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("dropCollection requiere nombre de colección")))
	}

	name := call.Argument(0).String()
	err := r.db.Collection(name).Drop(r.ctx)
	if err != nil {
		if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 26 {
			return goja.Undefined()
		}
		panic(r.vm.NewGoError(err))
	}

	return goja.Undefined()
}

func (r *scriptRunner) collection(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("collection requiere nombre")))
	}

	name := call.Argument(0).String()
	col := r.db.Collection(name)

	return r.buildCollectionObject(col)
}

func (r *scriptRunner) buildCollectionObject(col *mongo.Collection) *goja.Object {
	obj := r.vm.NewObject()
	_ = obj.Set("insertOne", r.wrapCollectionInsertOne(col))
	_ = obj.Set("insertMany", r.wrapCollectionInsertMany(col))
	_ = obj.Set("updateOne", r.wrapCollectionUpdateOne(col))
	_ = obj.Set("updateMany", r.wrapCollectionUpdateMany(col))
	_ = obj.Set("deleteOne", r.wrapCollectionDeleteOne(col))
	_ = obj.Set("deleteMany", r.wrapCollectionDeleteMany(col))
	_ = obj.Set("drop", r.wrapCollectionDrop(col))
	_ = obj.Set("createIndex", r.wrapCollectionCreateIndex(col))

	return obj
}

func (r *scriptRunner) runCommand(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("runCommand requiere documento")))
	}

	doc, err := r.valueToMap(call.Argument(0))
	if err != nil {
		panic(r.vm.NewGoError(err))
	}

	var result bson.M
	if err := r.db.RunCommand(r.ctx, mapToBsonD(doc)).Decode(&result); err != nil {
		panic(r.vm.NewGoError(err))
	}

	return r.vm.ToValue(result)
}

func (r *scriptRunner) wrapCollectionInsertOne(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(r.vm.NewGoError(fmt.Errorf("insertOne requiere documento")))
		}

		doc, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		if _, err := col.InsertOne(r.ctx, doc); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionInsertMany(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(r.vm.NewGoError(fmt.Errorf("insertMany requiere arreglo de documentos")))
		}

		docsVal, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		docs, ok := docsVal.([]interface{})
		if !ok {
			panic(r.vm.NewGoError(fmt.Errorf("insertMany espera arreglo de documentos")))
		}

		if _, err := col.InsertMany(r.ctx, docs); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionUpdateOne(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(r.vm.NewGoError(fmt.Errorf("updateOne requiere filtro y update")))
		}

		filter, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		update, err := r.exportValue(call.Argument(1))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		options := options.Update()
		if len(call.Arguments) > 2 {
			optsMap, err := r.valueToMap(call.Argument(2))
			if err != nil {
				panic(r.vm.NewGoError(err))
			}
			if optsMap != nil {
				if upsert, ok := optsMap["upsert"].(bool); ok {
					options = options.SetUpsert(upsert)
				}
			}
		}

		if _, err := col.UpdateOne(r.ctx, filter, update, options); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionUpdateMany(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(r.vm.NewGoError(fmt.Errorf("updateMany requiere filtro y update")))
		}

		filter, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		update, err := r.exportValue(call.Argument(1))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		options := options.Update()
		if len(call.Arguments) > 2 {
			optsMap, err := r.valueToMap(call.Argument(2))
			if err != nil {
				panic(r.vm.NewGoError(err))
			}
			if optsMap != nil {
				if upsert, ok := optsMap["upsert"].(bool); ok {
					options = options.SetUpsert(upsert)
				}
			}
		}

		if _, err := col.UpdateMany(r.ctx, filter, update, options); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionDeleteOne(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(r.vm.NewGoError(fmt.Errorf("deleteOne requiere filtro")))
		}

		filter, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		if _, err := col.DeleteOne(r.ctx, filter); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionDeleteMany(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(r.vm.NewGoError(fmt.Errorf("deleteMany requiere filtro")))
		}

		filter, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		if _, err := col.DeleteMany(r.ctx, filter); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionDrop(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if err := col.Drop(r.ctx); err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 26 {
				return goja.Undefined()
			}
			panic(r.vm.NewGoError(err))
		}
		return goja.Undefined()
	}
}

func (r *scriptRunner) wrapCollectionCreateIndex(col *mongo.Collection) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) == 0 {
			panic(r.vm.NewGoError(fmt.Errorf("createIndex requiere keys")))
		}

		keys, err := r.exportValue(call.Argument(0))
		if err != nil {
			panic(r.vm.NewGoError(err))
		}

		model := mongo.IndexModel{Keys: keys}
		if len(call.Arguments) > 1 {
			optsMap, err := r.valueToMap(call.Argument(1))
			if err != nil {
				panic(r.vm.NewGoError(err))
			}
			model.Options = buildIndexOptions(optsMap)
		}

		if _, err := col.Indexes().CreateOne(r.ctx, model); err != nil {
			panic(r.vm.NewGoError(err))
		}

		return goja.Undefined()
	}
}

func (r *scriptRunner) exportValue(value goja.Value) (interface{}, error) {
	var target interface{}
	if err := r.vm.ExportTo(value, &target); err != nil {
		return nil, err
	}
	return target, nil
}

func (r *scriptRunner) valueToMap(value goja.Value) (map[string]interface{}, error) {
	if goja.IsUndefined(value) || goja.IsNull(value) {
		return nil, nil
	}

	exported, err := r.exportValue(value)
	if err != nil {
		return nil, err
	}
	if exported == nil {
		return nil, nil
	}

	m, ok := exported.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("se esperaba objeto, obtuvo %T", exported)
	}
	return m, nil
}

type mongoDBObject struct {
	runner *scriptRunner
	props  map[string]goja.Value
}

func (d *mongoDBObject) Get(key string) goja.Value {
	if val, ok := d.props[key]; ok {
		return val
	}

	return d.runner.buildCollectionObject(d.runner.db.Collection(key))
}

func (d *mongoDBObject) Set(key string, val goja.Value) bool {
	d.props[key] = val
	return true
}

func (d *mongoDBObject) Has(key string) bool {
	if _, ok := d.props[key]; ok {
		return true
	}
	return true
}

func (d *mongoDBObject) Delete(key string) bool {
	delete(d.props, key)
	return true
}

func (d *mongoDBObject) Keys() []string {
	keys := make([]string, 0, len(d.props))
	for k := range d.props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func mapToBsonD(m map[string]interface{}) bson.D {
	if len(m) == 0 {
		return bson.D{}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	d := make(bson.D, 0, len(m))
	for _, k := range keys {
		d = append(d, bson.E{Key: k, Value: m[k]})
	}
	return d
}

func buildIndexOptions(opts map[string]interface{}) *options.IndexOptions {
	if len(opts) == 0 {
		return nil
	}

	indexOpts := options.Index()

	if name, ok := opts["name"].(string); ok {
		indexOpts = indexOpts.SetName(name)
	}
	if unique, ok := opts["unique"].(bool); ok {
		indexOpts = indexOpts.SetUnique(unique)
	}
	if sparse, ok := opts["sparse"].(bool); ok {
		indexOpts = indexOpts.SetSparse(sparse)
	}
	if expire, ok := opts["expireAfterSeconds"].(float64); ok {
		indexOpts = indexOpts.SetExpireAfterSeconds(int32(expire))
	} else if expireInt, ok := opts["expireAfterSeconds"].(int64); ok {
		indexOpts = indexOpts.SetExpireAfterSeconds(int32(expireInt))
	}

	return indexOpts
}

func (r *scriptRunner) isoDate(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("ISODate requiere string")))
	}

	input := call.Argument(0).String()
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		panic(r.vm.NewGoError(err))
	}

	return r.vm.ToValue(t)
}

func (r *scriptRunner) objectID(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("ObjectId requiere string")))
	}

	input := call.Argument(0).String()
	oid, err := primitive.ObjectIDFromHex(input)
	if err != nil {
		panic(r.vm.NewGoError(err))
	}

	return r.vm.ToValue(oid)
}

func (r *scriptRunner) print(call goja.FunctionCall) goja.Value {
	parts := make([]interface{}, 0, len(call.Arguments))
	for _, arg := range call.Arguments {
		parts = append(parts, arg.Export())
	}

	fmt.Println(parts...)
	return goja.Undefined()
}
