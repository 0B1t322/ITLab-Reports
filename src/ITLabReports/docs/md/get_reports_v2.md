# Description 

**Return a list of reports**

## Params
### Query
1. `offset` - the offset of the searching should be greater or equal 0, else it will be ignore
2. `limit` - the limit of getted reports, should be greater or equal 1, else it will be ignore
3. `sortBy` - the sorting of reports, should be in format `field_1:ordering`, where `ordering` can be: "asc", "ASC", "1" or "desc", "DESC", "-1". Can be passed more than one `sortBy`.
    3.1. `name:ordering` - sort by name
    3.2. `date:ordering` - sort by date
4. `dateBegin` - find reports whose date greater or equal then passed date. Date should be in format RFC3339 for example `dateBegin=2019-10-12T07:20:50.52Z`
5. `dateEnd` - find reports whose date lower or equal then passed date. Date should be in format RFC3339 for example `dateEnd=2019-10-12T07:20:50.52Z`
6. `match` - passed match params into filtering. match should be in format `field_1:value_1`, else it will be ignore. Can be passed more than one `match`.
    6.1. `name:value` - find reports whose name is like passed `value`
    6.2. `date:value` - find reports whose date is equal passed `value`, date should be in format RFC3339, for example `match=date:2019-10-12T07:20:50.52Z`
    6.3. `assignees.implementer:value` - find reports whose implementer is equal passed `value`, can work only if you are admin, else it will be unset
    6.4. `assignees.reporter:value` - find reports whose reporter is equal passed `value`, can work only if you are admin, else it will be unset
7. `onlyApproved` - filter reports on their approve state. If true return only reports that approved.
## Responce fields
1. `count` - the current count of elements in `items` field.
2. `items` - the getted reports.
3. `hasMore` - indecate that you can get new reports by increase `offset` param.
4. `limit` - the passed `limit` param.
5. `offset` - the passed `offset` param.
6. `totalResult` - how mush reports you can get by passing this request with this filtering params.