# Description 

**Create draft**

Create a draft with specifed `text` and `name`. `implementer` is optional query parameter, if not set in default - it's you.

## Params
### Query
1. `implementer` - id of user, that implements this draft

## Errors
### 400 - Bad Request
1. If `text` or `name` is empty
