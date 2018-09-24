import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {filter, get, isEmpty, map, values} from 'lodash'
import {formatDate} from '../../types/date'
import {requestReadIssues} from  '../../redux/actions/issues'
import './index.css'


const mapStateToProps = (state, ownProps) => ({
  issues: values(filter(get(state, `issuesByVolume`, {}), (v, k) => k !== 'latest')),
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

    return (
      <div className='main'>
        <div className='issues-summaries-container'>
          <div className='issues-header'>volumes</div>
          {
            map(
              issues,
              (issue, idx) => (
                <IssueRow issue={issue} key={idx}/>
              ),
              [],
            )
          }
        </div>
      </div>
    )
  }
}

export const Issues = connect(mapStateToProps, actions)(issues)

export class IssueRow extends Component {
  render() {
    const {
      volume,
      title,
      date,
      description,
    } = this.props.issue

    return (
      <Link to={`/issue/${volume}`}>
        <div className='issue-row-container'>
          <div className='issue-row'>
            <div className='issue-row-title-line'>
              <span>
                <span className='issue-row-volume-item'>{`Vol. ${volume}`}</span>
                <span className='issue-row-title-item'>{title}</span>
              </span>
              <span>{formatDate(date)}</span>
            </div>
            <div className='issue-row-description'>
              {description}
            </div>
          </div>
        </div>
      </Link>
    )
  }
}
