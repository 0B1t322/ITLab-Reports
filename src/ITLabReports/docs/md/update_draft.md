# Description 

**Update draft**

Update in draft only fields specifeid on body

You can update nothing but it's important to send empty json object:
```json
{

}
```

## Params
### Path
1. `id` - id of `draft` report

## Errors
### 404 - Not Found
1. If not found `draft` with `id` param.
### 400 - Bad Request
1. If `id` param is not valid.
2. Try to set `text` or `name` or `implementer` to empty.
### 403 - Forbiden
1. If you are not owner of this `draft`
