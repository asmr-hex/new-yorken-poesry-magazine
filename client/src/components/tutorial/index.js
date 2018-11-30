import React from 'react'
import {Mainframe} from '../ascii/mainframe'
import '../home/index.css'


export const Tutorial = props => {
  const pageStyle = {
    display: 'flex',
    flexDirection: 'column',
    width: '100%',
    margin: '2em 25% 5em 25%',
  }
  const headerStyle = {
    color: '#ffb2e4',
    fontSize: '4em',
    textShadow: '4px 4px #affbff',
  }
  const contentStyle = {
    textAlign: 'justify',
    textJustify: 'inter-word',
    lineHeight: '1.5em',
    fontSize: '2.3em',
    fontWeight: 'bold',
    color: '#19ecff',
    textShadow: '2px 2px #ffb2e4',
  }
  const comicStyle = {
    lineHeight: 'normal',
    fontWeight: 'normal',
    textShadow: 'none',
  }
  
  return (
    <div className={'main'}>
      <div style={pageStyle}>
        <div>
          <p className={'tutorial-header'} style={headerStyle}>tutorial</p>
          <div className={'tutorial-content'} style={contentStyle}>
            <div style={comicStyle}>
              <Mainframe size={'medium'}/>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
