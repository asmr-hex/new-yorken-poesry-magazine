import React, { Component } from 'react';
import {connect} from 'react-redux'
import {showMenu} from '../../redux/actions/ui'
import './title.css';

const actions = {
  showMenu,
}

class title extends Component {
  render() {
    return (
      <div onClick={() => this.props.showMenu()} className="header-title">
        New Yorken Poesry
      </div>
    )
  }
}

export const Title = connect(() => ({}), actions)(title)
