package pagelet

type Filter func(c *Controller)

// Filters is the default set of global filters.
// It may be set by the application on initialization.
var Filters = []Filter{
    RouterFilter,  // Use the routing table to select the right Action.
    ParamsFilter,  // Parse parameters into Controller.Params.
    SessionFilter, // Restore and write the session cookie.
    I18nFilter,    // Resolve the requested language.
    ActionInvoker, // Invoke the action.
}
