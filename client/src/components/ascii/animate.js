import React, {Component} from 'react';
import {map, range, reduce} from 'lodash'
import './glitch.css'


// this component is used to wrap ascii art sprites and apply animations
export class Animation extends Component {  
  constructor(props) {
    super(props)
    
    const defaultAnimationProps = {
      speed: 1,  // in seconds
      frames: ['Default A', 'Default B', 'Default C'],
      sequence: null,
      repeat: true,
      size: 'large',  // xx-small, x-small, small, ..., large, x-large, xx-large
      glitchEnabled: true,
      style: {},
      mainStyle: {
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

    // jeez....thanksgiving night...i ate too much tofurky
    // and i created this monstrosity... sorry earth.
    // welp.
    // obviously this...thing...returns the frame with the
    // largest rectangular volume. this is used to stabalize
    // the animations and is rendered with opacity 0.
    this.maxFrame = this.animation.frames[reduce(
      this.animation.frames,
      (acc, frame, idx) => {
        const length = reduce(
          frame.toString().split(/\r?\n/),
          (l, f) => l < f.length ? f.length : l,
          0,
        )
        const height = frame.toString().split(/\r?\n/).length
        
        return acc[1] < length*height ? [idx, length*height] : acc
      },
      [0, 0],
    )[0]]
    
  }

  componentDidMount() {
    if (this.animation.frames.length !== 1) {
      this.renderFrame()
    }
  }
  
  renderFrame() {
    this.setState({
      idx: (this.state.idx+1) % this.animation.sequence.length,
    })

    if (this.state.idx === this.animation.sequence.length-1 && !this.animation.repeat) {
      return
    }
    
    setTimeout(this.renderFrame.bind(this), this.animation.speed * 1000) 
  }

  render() {
    const containerStyle = {
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      textAlign: 'left',
      fontSize: this.animation.size,
      ...this.animation.style,
    }

    const layerStyles = [
      this.animation.mainStyle,
      this.animation.middleStyle,
      this.animation.bottomStyle
    ]

    return (
      <div className={this.animation.className} style={containerStyle}>
        <div style={{position: 'relative', alignSelf: 'flex-start'}}>
          <pre style={{position: 'relative', float: 'left', opacity: 0, margin: '0px'}}>
            {this.maxFrame}
          </pre>
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
