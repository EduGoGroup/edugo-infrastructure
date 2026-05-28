import os
import re

directory = "postgres/entities"
files = [
    f for f in os.listdir(directory) if f.endswith(".go") and not f.endswith("_test.go")
]


def process_file(filepath):
    with open(filepath, "r") as file:
        lines = file.readlines()

    out_lines = []
    in_struct = False
    for line in lines:
        if line.startswith("type ") and " struct {" in line:
            in_struct = True
            out_lines.append(line)
            continue
        if in_struct and line.strip() == "}":
            in_struct = False
            out_lines.append(line)
            continue

        if in_struct:
            # Match field line like: Name string `gorm:"not null" json:"name"`
            match = re.match(r"^(\s+)(\w+)(\s+)([\w\.\*\[\]]+)(\s+)`(.+)`\s*$", line)
            if match:
                spaces1, name, spaces2, typ, spaces3, tags = match.groups()

                if "validate:" in tags:
                    out_lines.append(line)
                    continue

                is_pointer = typ.startswith("*")
                is_slice = typ.startswith("[]")
                base_type = typ.replace("*", "").replace("[]", "")

                is_gorm_relation = False

                # Check for standard types vs relationships
                standard_types = {
                    "string",
                    "int",
                    "int8",
                    "int16",
                    "int32",
                    "int64",
                    "float32",
                    "float64",
                    "bool",
                    "time.Time",
                    "uuid.UUID",
                    "json.RawMessage",
                    "gorm.DeletedAt",
                    "byte",
                }

                if base_type not in standard_types:
                    is_gorm_relation = True

                # In GORM, associations can also be standard types if we made a mistake,
                # but typically they are pointers to other entities or slices of other entities.
                # E.g., `*School`, `[]UserRole`
                if is_pointer and base_type not in standard_types:
                    is_gorm_relation = True
                if is_slice and base_type not in {"byte"}:
                    is_gorm_relation = True

                validate_rules = []

                if is_gorm_relation:
                    validate_rules.append("-")
                elif base_type == "gorm.DeletedAt" or name in (
                    "CreatedAt",
                    "UpdatedAt",
                    "DeletedAt",
                ):
                    validate_rules.append("-")
                elif (
                    base_type == "bool"
                    or base_type == "json.RawMessage"
                    or base_type == "time.Time"
                ):
                    # Do not enforce required on bools, json or times
                    pass
                else:
                    is_required = (
                        "not null" in tags or "primaryKey" in tags
                    ) and not is_pointer

                    if is_pointer:
                        validate_rules.append("omitempty")
                    elif is_required:
                        validate_rules.append("required")
                    else:
                        validate_rules.append("omitempty")

                    if base_type == "uuid.UUID":
                        validate_rules.append("uuid")
                    elif "Email" in name:
                        validate_rules.append("email")
                    elif "Url" in name or "URL" in name:
                        validate_rules.append("url")
                    elif base_type == "string":
                        # Check enum
                        check_match = re.search(r" IN \((.*?)\)", tags, re.IGNORECASE)
                        if check_match:
                            vals = [
                                v.strip().strip("'").strip('"')
                                for v in check_match.group(1).split(",")
                            ]
                            validate_rules.append(f"oneof={' '.join(vals)}")
                        else:
                            if "Password" in name:
                                validate_rules.append("min=8")
                            else:
                                size_match = re.search(r"size:(\d+)", tags)
                                max_size = size_match.group(1) if size_match else "255"

                                if is_required:
                                    if (
                                        "Code" in name
                                        or "Country" in name
                                        or "City" in name
                                        or "Address" in name
                                        or "Phone" in name
                                        or "Name" in name
                                        or "Title" in name
                                        or "Description" in name
                                        or "FirstName" in name
                                        or "LastName" in name
                                    ):
                                        validate_rules.append("min=2")
                                        if max_size != "255":
                                            validate_rules.append(f"max={max_size}")
                                        else:
                                            validate_rules.append("max=255")
                                    else:
                                        if max_size != "255":
                                            validate_rules.append(f"max={max_size}")

                if validate_rules:
                    new_tags = tags + f' validate:"{",".join(validate_rules)}"'
                    out_lines.append(
                        f"{spaces1}{name}{spaces2}{typ}{spaces3}`{new_tags}`\n"
                    )
                else:
                    out_lines.append(line)
            else:
                out_lines.append(line)
        else:
            out_lines.append(line)

    with open(filepath, "w") as file:
        file.writelines(out_lines)


for filename in files:
    filepath = os.path.join(directory, filename)
    process_file(filepath)
print("Processed files successfully.")
