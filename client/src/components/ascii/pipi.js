import React from 'react';
import {get} from 'lodash'
import {Animation} from './animate'


export const PipiSauvage = (props) => {
  const TALKING = 'talking'
  const SWOONING = 'swooning'
  const WHISTLING = 'whistling'
  
  const defaults = {
    action: TALKING
  }

  const config = {
    ...defaults,
    ...props,
  }

  const actionFrames = {
    [TALKING]: getTalkingFrames(),
    [SWOONING]: getSwooningFrames(),
    [WHISTLING]: getWhistlingFrames(),
  }

  return (
    <Animation
      frames={get(actionFrames, config.action)}
      sequence={[0,1,2,1,2,1]}
      size={'x-large'}
      speed={1}
      {...config}
      />
  )
}

const getTalkingFrames = () =>
  ([String.raw`
     _____________________ 
    /383838383838338383838\ 
    |83/ ------------- \83| 
    |99||             ||99| 
    |83||   >◠    ◠<  ||98| 
    |83||             ||98|   
    |83||  0   )▾( 0  ||98|  
    |83||             ||98|   
    |83\ ------------- /98|    
     \8989898989898989898/     
          /888888888\          
    ========================   
   / 1 2 3 4 5 6 7 8 9 0 -  \\  
  / q w e r t y u i o p [ ]  \\
 / a s d f g h j k l ; ' entr \\ 
|  z x c v b n m , . ?  shift  ||
 -------------------------------`,
    String.raw`
     _____________________ 
    /383838383838338383838\ 
    |83/ ------------- \83| 
    |99||             ||99| 
    |83||   >●    ●<  ||98| 
    |83||             ||98|   
    |83||  0   )●( 0  ||98|  
    |83||             ||98|   
    |83\ ------------- /98|    
     \8989898989898989898/     
          /888888888\          
    ========================   
   / 1 2 3 4 5 6 7 8 9 0 -  \\  
  / q w e r t y u i o p [ ]  \\
 / a s d f g h j k l ; ' entr \\ 
|  z x c v b n m , . ?  shift  ||
 -------------------------------`,
    String.raw`
     _____________________ 
    /383838383838338383838\ 
    |83/ ------------- \83| 
    |99||             ||99| 
    |83||   >●    ●<  ||98| 
    |83||             ||98|   
    |83||  0   )▾( 0  ||98|  
    |83||             ||98|   
    |83\ ------------- /98|    
     \8989898989898989898/     
          /888888888\          
    ========================   
   / 1 2 3 4 5 6 7 8 9 0 -  \\  
  / q w e r t y u i o p [ ]  \\
 / a s d f g h j k l ; ' entr \\ 
|  z x c v b n m , . ?  shift  ||
 -------------------------------`,
   ])

const getWhistlingFrames = () =>
      ([String.raw`
     _____________________ 
    /383838383838338383838\ 
    |83/ ------------- \83| 
    |99||             ||99| 
    |83||  > ◠    ◠ < ||98| 
    |83||             ||98|   
    |83||  0   )o( 0  ||98|  
    |83||             ||98|   
    |83\ ------------- /98|    
     \8989898989898989898/     
          /888888888\          
    ========================   
   / 1 2 3 4 5 6 7 8 9 0 -  \\  
  / q w e r t y u i o p [ ]  \\
 / a s d f g h j k l ; ' entr \\ 
|  z x c v b n m , . ?  shift  ||
 -------------------------------
`])
  
const getSwooningFrames = () => 
      ([String.raw`
     _____________________ 
    /383838383838338383838\ 
    |83/ ------------- \83| 
    |99||             ||99| 
    |83||   ❤     ❤   ||98| 
    |83||  0   ◡   0  ||98|   
    |83||             ||98|  
    |83||             ||98|   
    |83\ ------------- /98|    
     \8989898989898989898/     
          /888888888\          
    ========================   
   / 1 2 3 4 5 6 7 8 9 0 -  \\  
  / q w e r t y u i o p [ ]  \\
 / a s d f g h j k l ; ' entr \\ 
|  z x c v b n m , . ?  shift  ||
 ------------------------------- 
`])
