# Description 

**Return a draft**
## Params
### Path
1. `id` - id of `draft` report

## Errors
### 404 - Not Found
1. If not found `draft` with `id` param.
### 400 - Bad Request
1. If `id` param is not valid.
### 403 - Forbiden
1. If you are not owner of this `draft`
