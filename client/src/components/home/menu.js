import React, { Component } from 'react';
import {connect} from 'react-redux'
import {showTitle} from '../../redux/actions/ui'
import './index.css';

const actions = {
  showTitle,
}

class menu extends Component {
  render() {
    return (
      <div className='home-menu'>
        <div>?</div>
        <div onClick={() => this.props.showTitle()} >~</div>
        <div>✐</div>
      </div>
    )
  }
}

export const Menu = connect(() => {}, actions)(menu)
