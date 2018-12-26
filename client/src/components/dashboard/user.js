import React, { Component } from 'react';
import {connect} from 'react-redux'
import {Chip} from '../ascii/chip'
import {SpeechBubble} from '../ascii/speech'


class userMgmt extends Component {
  render() {
    const {
      user,
    } = this.props

    const hellos = [
      'hiya', 'yo', 'sup',
      'qué onda', 'coucou',
      'おはよう', '哈罗', 'Приветик',
      'شكو ماكو', 'नमस्ते',
    ]
    
    const greeting = String.raw`${hellos[Math.floor(Math.random()*hellos.length)]} 
${user.username}!`

    return (
      <div className='profile-user-details-container'>
        <div className='profile-user-greeting'>
          <Chip
            mainStyle={{
              color: '#57c3ff',
            }}
            middleStyle={{
              color: '#42fffd',
            }}
            bottomStyle={{
              color: '#63fc45',
            }}/>
          <SpeechBubble
            style={{marginLeft: '2em'}}
            text={greeting}
            mainStyle={{
              color: '#5be6ff',
            }}
            middleStyle={{
              color: '#ffff5b',
            }}
            bottomStyle={{
              color: '#5bffc2',
            }}
            size={'40px'}
            speed={0.1}
            repeat={false}/>
        </div>
        <table className='profile-user-details'>
          <tr>
            <td className='profile-user-details-field'>email</td>
            <td className='profile-user-details-value'>{user.email}</td>
          </tr>
          <tr>
            <td className='profile-user-details-field'>password</td>
            <td className='profile-user-details-value'>**********</td>
          </tr>
        </table>
        <div className='profile-user-details-divider'></div>
      </div>
    )
  }
}

export const UserMgmt = connect(() => ({}), {})(userMgmt)
