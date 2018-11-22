import React, {Component} from 'react';
import {map, range} from 'lodash'
import './glitch.css'


// this component is used to wrap ascii art sprites and apply animations
export class Animation extends Component {  
  constructor(props) {
    super(props)
    
    const defaultAnimationProps = {
      speed: 1,  // in seconds
      frames: ['Default A', 'Default B', 'Default C'],
      sequence: null,
      size: 'large',  // xx-small, x-small, small, ..., large, x-large, xx-large
      glitchEnabled: true,
      style: {
        color: '#e58de8',
        opacity: 1,
      },
      middleStyle: {
        color: '#4ff9ff',
        opacity: 0.9,
        top: 1,
      },
      bottomStyle: {
        color: '#e58de8',
        opacity: 1,
        top: -1,
      },
    }
    
    this.animation = {
      ...defaultAnimationProps,
      ...this.props,
    }

    if (this.animation.sequence === null) {
      this.animation.sequence = range(0, this.animation.frames.length) 
    }

    this.state = { idx: 0 }
    
    if (this.animation.frames.length != 1) {
      this.renderFrame()
    }
    
  }

  renderFrame() {
    this.setState({
      idx: (this.state.idx+1) % this.animation.sequence.length,
    })
    setTimeout(this.renderFrame.bind(this), this.animation.speed * 1000)
  }

  render() {
    const containerStyle = {
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      textAlign: 'left',
      fontSize: this.animation.size,
    }

    const layerStyles = [this.animation.style, this.animation.middleStyle, this.animation.bottomStyle]

    return (
      <div style={containerStyle}>
        <div style={{position: 'relative', alignSelf: 'flex-start'}}>
          {
            map(
              range(0, this.animation.glitchEnabled ? layerStyles.length : 1),
              idx => (
                <pre className={`glitch-txt-${idx}`} style={layerStyles[idx]} key={idx}>
                  {this.animation.frames[this.animation.sequence[this.state.idx]]}
                </pre>                
              )
            )
          }
        </div>
      </div>
    )    
  }
}
