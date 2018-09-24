import React, {Component} from 'react'
import {Link} from 'react-router-dom'
import '../issues/index.css' // TODO (cw|923.2018) move the relevant parts to a css file here...


export class Poem extends Component {
  render() {
    const {
      title,
      author,
      content,
    } = this.props.poem

    return (
      <div className='issue-poem' id={this.props.elemId}>
        <div className='issue-poem-header'>
          <span className='issue-poem-title'>{title}</span>
          <span className='issue-poem-subheader'>
            <span>by</span>
            <Link to={`/poet/${author.id}`} className='text-link'>
              <span className='issue-poem-author'>{author.name}</span>
            </Link>
          </span>
        </div>
        <div className='issue-poem-content'>
          {content}
        </div>
      </div>
    )
  }
}
