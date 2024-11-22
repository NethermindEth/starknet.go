package typed

import (
	"fmt"
	"slices"
	"strings"
)

// TODO: create description
// Its just a basic validation, DO NOT RELLY FULLY ON IT, VERIFY YOURSELF
func (typedData *TypedData) ValidateTypedData() (isValid bool, err error) {
	//types
	err = validateTypes(typedData.Types, typedData.PrimaryType)
	if err != nil {
		return false, err
	}

	//primary type
	_, ok := typedData.Types[typedData.PrimaryType]
	if !ok {
		return false, fmt.Errorf("primary type not found in 'types'")
	}

	//domain
	domainType, ok := typedData.Types[RevisionV0.Domain()]
	if !ok {
		domainType = typedData.Types[RevisionV1.Domain()]
	}

	err = validateDomain(typedData.Domain, domainType)
	if err != nil {
		return false, err
	}

	//TODO: validate message

	return true, nil
}

func validateTypes(types map[string]TypeDefinition, primaryType string) error {
	var customTypes []string
	isEnum := make(map[string]bool)

	switch {
	case len(types) == 0:
		return fmt.Errorf("'types' cannot be empty")
	case len(types) == 1:
		return fmt.Errorf("'types' should have at least 2 fields (domain separator and primary type)")
	case len(types) >= 2:
		for key, typeDef := range types {
			if err := validateTypeName(key); err != nil {
				return err
			}

			//search for enums
			for _, param := range typeDef.Parameters {
				if param.Type == "enum" || param.Type == "enum*" {
					isEnum[param.Contains] = true
				}
			}

			customTypes = append(customTypes, key)
		}

		_, ok1 := types[RevisionV0.Domain()]
		_, ok2 := types[RevisionV1.Domain()]
		if !(ok1 || ok2) {
			return fmt.Errorf("missing domain separator type")
		}

		isTypeUsed := make(map[string]bool)
		//validate params
		for key, typeDef := range types {
			for _, param := range typeDef.Parameters {
				if isEnum[key] {
					if !isValidEnumType(param.Type, customTypes) {
						return fmt.Errorf("invalid enum param type of '%s'", param.Name)
					}
				} else if !isStandardType(param.Type) {
					if !isCustomType(param.Type, customTypes) {
						return fmt.Errorf("invalid param type")
					}
					isTypeUsed[param.Type] = true
				}

				paramTypeName, _ := strings.CutSuffix(param.Type, "*")
				if len(param.Contains) == 0 {
					if paramTypeName == "merkletree" || paramTypeName == "enum" {
						return fmt.Errorf("the parameter 'contains' needs to be specified at '%s'", param.Name)
					}
				} else {
					switch paramTypeName {
					case "merkletree":
						_, ok := types[param.Contains]
						if !ok {
							return fmt.Errorf("type '%s' stated at 'contains' not found in 'types' field", param.Contains)
						}

						isTypeUsed[param.Contains] = true
					case "enum":
						enumDef, ok := types[param.Contains]
						if !ok {
							return fmt.Errorf("type '%s' stated at 'contains' not found in 'types' field", param.Contains)
						}

						for _, enumParam := range enumDef.Parameters {
							if !isValidEnumType(enumParam.Type, customTypes) {
								return fmt.Errorf("invalid type '%s', all enum variants types must be enclosed in parenthesis", enumParam.Type)
							}
						}

						isTypeUsed[param.Contains] = true
					default:
						return fmt.Errorf("the type '%s' does not use 'contains'", paramTypeName)
					}

				}
			}
		}

		//check for 'dangling types'
		isTypeUsed[RevisionV0.Domain()] = true
		isTypeUsed[RevisionV1.Domain()] = true
		isTypeUsed[primaryType] = true

		for _, typeName := range customTypes {
			if !isTypeUsed[typeName] {
				return fmt.Errorf("all the types defined must be referenced by another type (no dangling types)")
			}
		}
	}

	return nil
}

// Checks if the provided type name is of a type defined in the 'types' field represented by a string slice, also validates arrays
func isCustomType(typeName string, customTypes []string) bool {
	//to validate arrays
	typeName, _ = strings.CutSuffix(typeName, "*")
	return slices.Contains(customTypes, typeName)
}

// Checks if the provided type name follows the rules of the enum type and if it is a standard or custom type; returns false otherwise
func isValidEnumType(typeName string, customTypes []string) bool {
	if strings.HasPrefix(typeName, "(") && strings.HasSuffix(typeName, ")") {
		typeName = typeName[1 : len(typeName)-1]
		enumTypes := strings.Split(typeName, ",")

		if enumTypes[0] == "" {
			return true
		}

		for _, typeName := range enumTypes {
			return isCustomType(typeName, customTypes) || isStandardType(typeName)
		}
	}

	return false
}

// ref: https://github.com/starknet-io/SNIPs/blob/5d5a42c654c27b377d8b7f90b453065fd19ec2eb/SNIPS/snip-12.md#type-identification
func validateTypeName(customTypeName string) error {
	switch {
	case customTypeName == "":
		return fmt.Errorf("no empty name")
	case slices.Contains(revision_0_basic_types, customTypeName) || slices.Contains(revision_1_basic_types, customTypeName):
		return fmt.Errorf("name can't match basic types like felt, ClassHash, timestamp, u128")
	case slices.Contains(revision_1_preset_types, customTypeName):
		return fmt.Errorf("name can't match preset types like TokenAmount, NftId, u256")
	case strings.HasSuffix(customTypeName, "*"):
		return fmt.Errorf("name can't end in *")
	case strings.HasPrefix(customTypeName, "(") && strings.HasSuffix(customTypeName, ")"):
		return fmt.Errorf("name can't be enclosed in parenthesis")
	case strings.Contains(customTypeName, ","):
		return fmt.Errorf("name can't contain the comma (,) character (since it is used as a delimiter in the enum type)")
	}

	return nil
}

func validateDomain(domain Domain, domainType TypeDefinition) error {
	if !(len(domainType.Parameters) >= 3 && len(domainType.Parameters) <= 4) {
		return fmt.Errorf("domain should only have 3 or 4 fields")
	}

	fieldNames := []string{"name", "version", "chainId", "revision"}

	switch {
	case domainType.Name == RevisionV0.domain:
		if domain.Revision != RevisionV0.Version() {
			return fmt.Errorf("invalid revision version: revision for '%s' should be '%d' but is '%d' ", domainType.Name, RevisionV0.Version(), domain.Revision)
		}

		for _, v := range domainType.Parameters {
			if !slices.Contains(fieldNames, v.Name) {
				return fmt.Errorf("invalid field name '%s'", v.Name)
			}
			if v.Type != "felt" {
				return fmt.Errorf("invalid field type '%s'. Should be 'felt'", v.Type)
			}
		}
	case domainType.Name == RevisionV1.domain:
		if domain.Revision != RevisionV1.Version() {
			return fmt.Errorf("invalid revision version: revision for '%s' should be '%d' but is '%d' ", domainType.Name, RevisionV1.Version(), domain.Revision)
		}

		for _, v := range domainType.Parameters {
			if !slices.Contains(fieldNames, v.Name) {
				return fmt.Errorf("invalid field name '%s'", v.Name)
			}
			if v.Type != "shortstring" {
				return fmt.Errorf("invalid field type '%s'. Should be 'shortstring'", v.Type)
			}
		}
	default:
		return fmt.Errorf("invalid Domain separator name")
	}

	if domain.Name == "" {
		return fmt.Errorf("domain name field is empty")
	}
	if domain.ChainId == "" {
		return fmt.Errorf("chainId field is empty")
	}
	if domain.Version == "" {
		return fmt.Errorf("version field is empty")
	}

	return nil
}
