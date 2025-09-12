# XSD to Proto Converter

Converts DDEX XSD schemas to Protocol Buffer definitions with XML tag annotations.

## XSD Primitive to Proto Mappings

Based on analysis of all DDEX schemas in `xsd/`, here are the XSD patterns we need to handle:

### 1. Basic XSD Types → Proto Types

| XSD Type | Proto Type | Example | Notes |
|----------|------------|---------|-------|
| `xs:string` | `string` | `<xs:element name="MessageId" type="xs:string"/>` | Most common type |
| `xs:integer` | `int32` | `<xs:element name="BitsPerSample" type="xs:integer"/>` | Whole numbers |
| `xs:boolean` | `bool` | `<xs:element name="IsProvidedInDelivery" type="xs:boolean"/>` | true/false |
| `xs:decimal` | `string` | `<xs:element name="Value" type="xs:decimal"/>` | Preserve precision |
| `xs:duration` | `string` | `<xs:element name="Duration" type="xs:duration"/>` | ISO 8601 duration (PT30S) |
| `xs:dateTime` | `string` | `<xs:element name="DateTime" type="xs:dateTime"/>` | ISO 8601 datetime |
| `xs:date` | `string` | `<xs:element name="Date" type="xs:date"/>` | ISO 8601 date (YYYY-MM-DD) |
| `xs:anyURI` | `string` | `<xs:element name="URL" type="xs:anyURI"/>` | URI/URL as string |

### 2. XSD Structure → Proto Structure

#### A. Complex Type with Sequence
```xml
<xs:complexType name="MessageHeader">
  <xs:sequence>
    <xs:element name="MessageId" type="xs:string"/>
    <xs:element name="MessageCreatedDateTime" type="xs:dateTime"/>
  </xs:sequence>
</xs:complexType>
```
**→ Proto:**
```protobuf
message MessageHeader {
  string message_id = 1 [(tagger.tags) = "xml:\"MessageId\""];
  string message_created_date_time = 2 [(tagger.tags) = "xml:\"MessageCreatedDateTime\""];
}
```

#### B. Choice → Oneof
```xml
<xs:choice>
  <xs:element name="ISRC" type="xs:string"/>
  <xs:element name="ISWC" type="xs:string"/>
  <xs:element name="ProprietaryId" type="ern:ProprietaryId"/>
</xs:choice>
```
**→ Proto:**
```protobuf
oneof identifier {
  string isrc = 1 [(tagger.tags) = "xml:\"ISRC\""];
  string iswc = 2 [(tagger.tags) = "xml:\"ISWC\""];
  ProprietaryId proprietary_id = 3 [(tagger.tags) = "xml:\"ProprietaryId\""];
}
```

#### C. Cardinality → Repeated
```xml
<!-- maxOccurs="unbounded" becomes repeated -->
<xs:element name="ReleaseAdmin" minOccurs="0" maxOccurs="unbounded" type="ern:ReleaseAdmin"/>
```
**→ Proto:**
```protobuf
repeated ReleaseAdmin release_admin = 1 [(tagger.tags) = "xml:\"ReleaseAdmin\""];
```

#### D. Attributes → Fields with attr tag
```xml
<xs:complexType name="SomeType">
  <xs:simpleContent>
    <xs:extension base="xs:string">
      <xs:attribute name="IsDefault" type="xs:boolean" use="required"/>
    </xs:extension>
  </xs:simpleContent>
</xs:complexType>
```
**→ Proto:**
```protobuf
message SomeType {
  string value = 1 [(tagger.tags) = "xml:\",chardata\""];
  bool is_default = 2 [(tagger.tags) = "xml:\"IsDefault,attr\""];
}
```

#### E. Simple Elements
```xml
<xs:element name="URL" minOccurs="0" maxOccurs="unbounded" type="xs:string"/>
```
**→ Proto:**
```protobuf
repeated string url = 1 [(tagger.tags) = "xml:\"URL\""];
```

### 3. Cardinality Rules

| XSD | Proto | Meaning |
|-----|-------|---------|
| `minOccurs="0"` | Optional field | Field can be omitted |
| `maxOccurs="unbounded"` | `repeated` | Array/slice of values |
| `minOccurs="1" maxOccurs="1"` | Regular field | Required field (default) |
| `use="required"` (attributes) | Regular field | Required attribute |
| `use="optional"` (attributes) | Optional field | Optional attribute |

### 4. Custom Types → Message References

```xml
<xs:element name="MessageHeader" type="ern:MessageHeader"/>
```
**→ Proto:**
```protobuf
MessageHeader message_header = 1 [(tagger.tags) = "xml:\"MessageHeader\""];
```

The namespace prefix (`ern:`) is stripped and the type becomes a message reference.

## Schema Statistics

- **ERN v4.3.2**: 206 complex types, 374 unbounded elements, 91 simpleContent types
- **MEAD v1.1**: 53 simpleContent types  
- **PIE v1.0**: 42 simpleContent types

## Implementation Priority

1. **complexType with sequence** - Most common pattern
2. **Basic type mapping** - Handle xs:string, xs:integer, etc.
3. **Cardinality** - Handle minOccurs/maxOccurs  
4. **Attributes** - Handle simpleContent with attributes
5. **Choice** - Handle oneof patterns
6. **Custom types** - Reference other messages

The goal: Generate `.proto` files that produce Go structs with identical XML marshaling to the current `ddex/` packages.