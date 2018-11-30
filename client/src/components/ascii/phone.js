import React from 'react';
import {get, map, range} from 'lodash'
import {Animation} from './animate'


export const Phone = (props) => {
  const TALKING = 'talking'
  
  const defaults = {
    action: TALKING
  }

  const config = {
    ...defaults,
    ...props,
  }

  const actionFrames = {
    [TALKING]: getTalkingFrames(),
  }

  return (
    <Animation
      frames={get(actionFrames, config.action)}
      size={'x-large'}
      speed={1}
      {...config}
      />
  )
}

const getTalkingFrames = () =>
      ([String.raw`
      _________
     /####-####\\
    ｜          ||
    ｜          ||
    ｜          ||
    ｜  ● ‿ ●   ||
    ｜          ||
    ｜          ||
    ｜          ||
    ｜          ||
     \####◯####//
      ▔▔▔▔▔▔▔▔▔`,
        String.raw`
      _________
     /####-####\\
    ｜          ||
    ｜          ||
    ｜          ||
    ｜  ^ ‿ ^   ||
    ｜          ||
    ｜          ||
    ｜          ||
    ｜          ||
     \####◯####//
      ▔▔▔▔▔▔▔▔▔`,
       ])
