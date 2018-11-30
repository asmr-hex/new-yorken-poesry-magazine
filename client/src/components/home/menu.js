import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'
import {showTitle} from '../../redux/actions/ui'
import {Animation} from '../ascii/animate'
import './menu.css';

const actions = {
  showTitle,
}

// the menu should allow the user to navigate the poetic universe
// home | ./about | ./volumes | ./poets | ./sign{in,up} (or user profile pic)
class menu extends Component {
  render() {
    const {
      loggedIn,
    } = this.props

    // TODO (cw|11.24.2018) refactor this stuff below. possible render all these
    // links in a loop so we don't have to write everything out.
    return (
      <div className='home-menu'>
        <Link to='/' className='header-menu-item'>
          <Animation
            frames={['> █', '>']}
            style={{fontSize: 'inherit'}}
            mainStyle={{color: '#ffbae5'}}
            bottomStyle={{color: '#f4f3a4'}}
            />
        </Link>
        <Link to='/about' className='header-menu-item'>
          <Animation
            frames={['about']}
            style={{fontSize: 'inherit', fontWeight: 800}}
            mainStyle={{color: '#ffbae5', opacity: 1}}
            bottomStyle={{color: '#f4f3a4', opacity: 1}}/>
        </Link>
        <Link to='/tutorial' className='header-menu-item'>
          <Animation
            frames={['tutorial']}
            style={{fontSize: 'inherit', fontWeight: 800}}
            mainStyle={{color: '#ffbae5', opacity: 1}}
            bottomStyle={{color: '#f4f3a4', opacity: 1}}/>
        </Link>
        <Link to='/issues' className='header-menu-item'>
          <Animation
            frames={['volumes']}
            style={{fontSize: 'inherit', fontWeight: 800}}
            mainStyle={{color: '#ffbae5', opacity: 1}}
            bottomStyle={{color: '#f4f3a4', opacity: 1}}/>
        </Link>
        <Link to='/poets' className='header-menu-item'>
          <Animation
            frames={['poets']}
            style={{fontSize: 'inherit', fontWeight: 800}}
            mainStyle={{color: '#ffbae5', opacity: 1}}
            bottomStyle={{color: '#f4f3a4', opacity: 1}}/>
        </Link>
        {
          loggedIn ?
            <Link to='/profile' className='header-menu-item'>
                @
            </Link>
            : <Link to='/login' className='header-menu-item'>
                <Animation
                    frames={['signup', 'signin']}
                    speed={2}
                    style={{fontSize: 'inherit', fontWeight: 800}}
                    mainStyle={{color: '#ffbae5', opacity: 1}}
                    bottomStyle={{color: '#f4f3a4', opacity: 1}}/>
              </Link>
        }
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
