// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpsrv

type Filter func(c *Controller)

// Filters is the default set of global filters.
// It may be set by the application on initialization.
var DefaultFilters = []Filter{
	RouterFilter,  // Use the routing table to select the right Action.
	ParamsFilter,  // Parse parameters into Controller.Params.
	SessionFilter, // Restore and write the session cookie.
	I18nFilter,    // Resolve the requested language.
	ActionInvoker, // Invoke the action.
}
