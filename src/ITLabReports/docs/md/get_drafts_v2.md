# Description 

**Return a list of drafts**

## Params
### Query
1. `offset` - the offset of the searching should be greater or equal 0, else it will be ignore
2. `limit` - the limit of getted reports, should be greater or equal 1, else it will be ignore

## Responce fields
1. `count` - the current count of elements in `items` field.
2. `items` - the results drafts.
3. `hasMore` - indicate that you can get new drafts by increase `offset` param.
4. `limit` - the passed `limit` param.
5. `offset` - the passed `offset` param.
6. `totalResult` - how mush drafts you can get by passing this request with this filtering params.