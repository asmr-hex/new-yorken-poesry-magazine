import React, { Component } from 'react';
import {Route, Switch} from 'react-router-dom'
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {map} from 'lodash'
import {showTitle} from '../../redux/actions/ui'
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
  contributors: getContributorsByIssueVolume(ownProps.volume, state),
  poems: getPoemsByIssueVolume(ownProps.volume, state),
})

const actions = {
  fetchIssue: requestReadIssue,
  showTitle,
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
    return (
      <div className='Issue'>
        <TOC issue={this.props.issue} poems={this.props.poems} showTitle={this.props.ShowTitle}/>
        {
          map(
            this.props.poems,
            (i, idx) => <Poem poem={i} key={idx}/>,
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
    const {
      title,
      date,
      description,
      volume,
    } = this.props.issue
    
    return (
      <div className='toc'>
        <h1>{title}</h1>
        <h3>Vol. {volume}</h3>
        <h5>{date}</h5>
        <h3>{description}</h3>
        <div>
          {
            map(
              this.props.poems,
              (poem, idx) => (
                <div key={idx}>
                  {poem.title}
                  ............
                  <Link to={`/poet/${poem.author.id}`} className='text-link'>
                    {poem.author.name}
                  </Link>
                </div>
              ),
              [],
            )
          }
        </div>
      </div>
    )
  }
}

// TODO (cw|9.22.2018) move this to component/poems directory
class Poem extends Component {
  render() {
    const {
      title,
      date,
      author,
      content,
      likes,
    } = this.props.poem
    return (
      <div>
        <h2>{title}</h2>
        <h2>{author.name}</h2>
        <h6>{content}</h6>
      </div>
    )
  }
}
