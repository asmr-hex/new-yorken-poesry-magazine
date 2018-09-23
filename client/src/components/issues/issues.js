import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {get, isEmpty, map, values} from 'lodash'
import {requestReadIssues} from  '../../redux/actions/issues'


const mapStateToProps = (state, ownProps) => ({
  issues: values(get(state, `issuesByVolume`, {})),
})

const actions = {
  fetchIssues: requestReadIssues,
}

export class issues extends Component {
  componentDidMount() {
    const {
      issues,
      fetchIssues,
    } = this.props

    // get issues if we haven't already
    if (isEmpty(issues)) {
      fetchIssues()
    }
  }
  
  render() {
    const {
      issues,
    } = this.props

    console.log(issues)
    
    return (
      <div>
        {
          map(
            issues,
            issue => (
              <div>{issue.title}</div> 
            ),
            [],
          )
        }
      </div>
    )
  }
}

export const Issues = connect(mapStateToProps, actions)(issues)
