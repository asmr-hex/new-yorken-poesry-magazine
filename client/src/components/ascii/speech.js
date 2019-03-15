import React from 'react'
import {range, reduce} from 'lodash'
import {Animation} from './animate'


// this component is for rendering animated speech bubbles which
// are displayed character by character
export const SpeechBubble = (props) => {
  const {text} = props

  // separate the text chars into individual and
  // cumulative frames!
  const frames = reduce(
    range(0, text.length),
    (acc, idx) => ([
      ...acc,
      idx === 0 ?
        text.charAt(idx)
        : `${acc[idx-1]}${text.charAt(idx)}`,
    ]),
    [],
  )

  return (
    <Animation frames={frames}
               {...props}
               />
  )

}
