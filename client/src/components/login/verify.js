import React, { Component } from 'react';
import {connect} from 'react-redux'
import {requestVerification} from '../../redux/actions/verify'
import queryString from 'query-string'


const mapStateToProps = (state, ownProps) => {
  // parse query string
  const {token, email} = queryString.parse(ownProps.location.search)
  
  return {
    token,
    email,
  }
}

const actions = {
  requestVerification,
}

class verify extends Component {
  constructor(props) {
    super(props)
  }

  componentDidMount() {
    // immediately send a request to verify given the token
    // and email address in the query parameters of the URL
    this.props.requestVerification(
      {
        token: this.props.token,
        email: this.props.email,
      },
      
      this.redirectUponVerification,
    )
  }

  // upon verification, we want to redirect the route to the dashboard
  redirectUponVerification = () => {
    const {history} = this.props

    history.push('/profile')
  }

  render() {
    return (
      <div>
        WE ARE DOING ARE BEST. HOLD ON TO YOUR HORSES.
      </div>
    )
  }
  
}

export const Verify = connect(mapStateToProps, actions)(verify)

