import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {map} from 'lodash'
import {requestReadIssue} from '../../redux/actions/issues'
import {getIssueByVolume} from '../../redux/selectors/issues'
import {
  getJudgesByIssueVolume,
  getContributorsByIssueVolume,
} from '../../redux/selectors/poets'
import {getPoemsByIssueVolume} from '../../redux/selectors/poems'


const mapStateToProps = (state, ownProps) => ({
  issue: getIssueByVolume(ownProps.volume, state),
  judges: getJudgesByIssueVolume(ownProps.volume, state),
  contributors: getContributorsByIssueVolume(ownProps, state),
  poems: getPoemsByIssueVolume(ownProps, state),
})

const actions = {
  fetchIssue: requestReadIssue,
}

// Issue represents the a single issue view including different modes
// such as table of contents, poems, and additional information (e.g.
// committee of judges and contributors)
export class issue extends Component {
  componentDidMount() {
    const {fetchIssue, volume} = this.props

    // fetch this issue!
    fetchIssue(volume)
  }
  
  render() {
    const poems = []
    return (
      <div className='Issue'>
        <TOC/>
        {
          map(
            poems,
            i => <Poem data={i}/>,
            [],
          )
        }
      </div>
    )
  }
}

export const Issue = connect(mapStateToProps, actions)(issue)

class TOC extends Component {
  render() {
    return (
      <div className='toc'>
        table of contents
      </div>
    )
  }
}

class Poem extends Component {
  render() {
    return (
      <div/>
    )
  }
}
