import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {showTitle} from '../../redux/actions/ui'
import './menu.css';

const actions = {
  showTitle,
}

class menu extends Component {
  render() {
    
    return (
      <div className='home-menu'>
        <Link to='/about' className='header-menu-item' onClick={() => this.props.showTitle()}>?</Link>
        <Link to='/about' className='header-menu-item' onClick={() => this.props.showTitle()}>#</Link>
        <Link to='/' className='header-menu-item' onClick={() => this.props.showTitle()}>~</Link>
        <Link to='/about' className='header-menu-item' onClick={() => this.props.showTitle()}>&</Link>
        <Link to='login' className='header-menu-item' onClick={() => this.props.showTitle()}>@</Link>
      </div>
    )
  }
}

/*
  did i find a bug in react router? basically, if i link to a route in component B
  and the route rendered occurs in component A and component B is a child of A, then
  the route doesn't updated unless i trigger some re-render in the parent component
  somehow. Need to look into this and maybe submit PR for bugfix.
 */

export const Menu = connect(() => ({}), actions)(menu)
