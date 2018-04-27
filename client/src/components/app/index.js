import React, { Component } from 'react';
import {Route, Switch, withRouter} from 'react-router-dom'
import {connect} from 'react-redux'
import {Home} from '../home'
import {Dashboard} from '../dashboard'


class app extends Component {
  render() {
    return (
      <Switch>
        {/* we *must* specify routes in a specific order so the matching works -__- */}
        <Route path='/dashboard' component={Dashboard}/>
        <Route path='/' component={Home}/>
      </Switch>
    )
  }
}

// Note (cw|4.26.2018): updating the url isn't updating the components for
// some reason, we need to explicitly pass the new location information of
// the url to this component. this is outline here
// https://github.com/ReactTraining/react-router/blob/master/packages/react-router/docs/guides/blocked-updates.md
// with a quick solution of using `withRouter`, though it is not the most
// efficient solution (see
// https://github.com/ReactTraining/react-router/pull/5552#issuecomment-331502281)
// which explains that this will update all components even if they shouldn't
// be re-rendered? In our case this seems permissible, but the most efficient
// solution would be to thread location through as a prop.
export const App = withRouter(connect(() => ({}), {})(app))
