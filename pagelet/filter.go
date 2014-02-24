package pagelet

type Filter func(c *Controller, filterChain []Filter)

// Filters is the default set of global filters.
// It may be set by the application on initialization.
var Filters = []Filter{
    RouterFilter,  // Use the routing table to select the right Action.
    I18nFilter,    // Resolve the requested language.
    ActionInvoker, // Invoke the action.
}
