# Description 

**Create a report**

Have `implementer` query param that Indicates who the report is about, if not specifed - by default it's you.

`name` field is optional so you can create report with two ways:

## Creating by only ```text``` field

To create by only text field you should to set in text field ```name``` of report and field ```text``` and separet it with ```@\n\t\n@```
For example next body will create report with ```name``` report_text and with ```text``` report text
```json
{
    "text": "report_text@\n\t\n@report_text"
}
```

## Creating by ```text``` and ```name``` fields

It's more **preferred** method just send next body:
```json
{
    "name": "report_name",
    "text": "report_text"
}
```

# Errors

## 400 Bad request
1. If `name` and `text` fields is empty
2. If `name` is not empty, `text` is empty

