package pagelet

type Filter func(c *Controller, filterChain []Filter)

// Filters is the default set of global filters.
// It may be set by the application on initialization.
var Filters = []Filter{
    RouterFilter, // Use the routing table to select the right Action.
    //ParamsFilter,            // Parse parameters into Controller.Params.
    ActionInvoker, // Invoke the action.
}
