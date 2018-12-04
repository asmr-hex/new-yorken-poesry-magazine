import React from 'react'
import Highlight from 'react-highlight'
import {Mainframe} from '../ascii/mainframe'
import {Chip} from '../ascii/chip'
import {Floppy} from '../ascii/floppy'
import {Phone} from '../ascii/phone'
import {SpeechBubble} from '../ascii/speech'
import '../home/index.css'
import '../app/highlight.css'


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
  const subheaderStyle = {
    color: '#ffb2e4',
    fontSize: '3em',
    textShadow: '4px 4px #affbff',
  }
  const contentStyle = {
    textAlign: 'justify',
    textJustify: 'inter-word',
    lineHeight: '1.3em',
    fontSize: '2.5em',
    fontWeight: 'bold',
    color: '#19ecff',
    textShadow: '2px 2px #ffb2e4',
    marginTop: '1em',
    marginBottom: '3em',
  }
  const keywordStyle = {
    fontStyle: 'italic',
    color: '#75b0ff',
  }
  const codeStyle = {
    fontWeight: 'normal',
    fontSize: '0.8em',
    backgroundColor: '#f4f4f4',
    color: '#6272a4',
    textShadow: 'none',
    // wordSpacing: '2em',
  }
  const comicStyle = {
    display: 'flex',
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
            <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
              <Mainframe
                style={comicStyle}
                glitchEnabled={false}
                // mainStyle={{color: '#66ddb7'}}
                mainStyle={{color:'#ff75b3', fontWeight: 'bold', textShadow: '6px 6px #19ecff'}}
                size={'small'}/>
              <SpeechBubble
                speed={0.1}
                repeat={false}
                glitchEnabled={false}
                style={{fontSize: '1em'}}
                mainStyle={{color: '#19ecff', textShadow: '3px 3px #ffb2e4'}}
                text={String.raw`making ur 0wn 
algopoetic pa1 
is 3asy and fun!`}/>
            </div>
          </div>
          <div style={subheaderStyle}>intr0duction</div>
          <div className={'tutorial-content'} style={contentStyle}>
            {
              // yes! it is a computer program you can write to write
              // poetry! BRIEF NYPM INTRO (POINT TO ABOUT PAGE FOR DEETZ) there is a rich history of computer generated
              // poetry which we will not discuss here, but plz feel
              // free to check out some of the <a href='#resources'>resources</a> provided.
              // <br/>
              //   <br/>
                  
            }
            in this tutorial, we will show how you can
            write your own algorithmic poet which can be submitted to
            the <em>New Yorken Poesry Magazine</em>. we assume a beginner's
            understanding of the Python programming language, in particular
            <ul style={{fontSize: '0.7em'}}>
              <li>lists/dicts</li>
              <li>functions</li>
              <li>parsing command-line arguments</li>
            </ul>
            and basic familiarity with running programs from the command-line.
            
            
          </div>
          <div style={subheaderStyle}>what's an algorithmic poet?</div>
          <div className={'tutorial-content'} style={contentStyle}>
            an <em>algorithmic poet</em> is a computer program you create that
            can <em style={keywordStyle}>write poems</em>, <em style={keywordStyle}>critique poems</em>,
      & <em style={keywordStyle}>learn how to write better poems</em>! your program
    should be able to handle each of the following command-line arguments in order to perform these tasks.
      <div style={{display: 'flex', justifyContent: 'center'}}>
      <ul style={{width: '70%', listStyleType: 'none'}}>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--write</span> generate and print a poem to standard out (stdout).
      </li>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--critique POEM</span> critique a given poem, POEM, and print a score between
    0 (worst) and 1 (best).
      </li>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--study POEM</span> read a given poem, POEM, and <em>optionally</em> use
      it to modify how your program writes poetry. this task should not print any output to stdout.
      </li>
      </ul>
      </div>
      for example, if you write a poet in Python called <span style={codeStyle}>poet_bot.py</span>,
    then you can ask it to perform each task by running it on the command line like this,
      <Highlight className="shell">
      {
        String.raw`$ python poet_bot.py --write  # write a poem
roses are red, violets are blue

$ python poet_bot.py --critique "this is a bad poem" # critique a poem
0.88

$ python poet_bot.py --study "this is a poem to study"  # learn how to write better
`
      }
      </Highlight>
      if these tasks seem kinda vague, don't worry! they are like that for a reason. as long as
      your poet satisfies the command-line <em>interfaces</em> described above, it can be implemented
      in whatever way you choose. this gives a poet designer (like you!)
      nearly limitless options for deciding how their poet will work!
      <br/><br/>
      for the remainder of this tutorial we will go over in detail a couple ways to implement each of these tasks! 
          </div>
          <div style={subheaderStyle}>your first poet</div>
          <div className={'tutorial-content'} style={contentStyle}>
            let's begin by making a simple poet which can perform the required tasks. this implementation
            will not be super interesting, but it will give us a foundation to build more complex poets.
            <br/><br/>
            we'll start by making a new empty Python file called <span style={codeStyle}>poet_bot.py</span>.
            you can name your own file differently, but this is what we will be calling our program
            in this tutorial.
            <br/><br/>
            great! since our program needs to be able to perform three distinct tasks,
            let's define three functions corresponding to each. in your favorite text editor or integrated
            developer environment (IDE), type out the following in your Python file:

            <Highlight className="python poet-body-code">
              {
                String.raw`# poet_bot.py

# this function generates and prints a poem.
def write_poem():
    print('roses are red, violets are blue')

# when given a poem, this function critiques it on
# a scale from 0-1 and prints the score.
def critique_poem(a_poem):
   print(0.88)

# when given a poem, this function can allow our
# program to potentially learn new styles or approaches to
# writing poetry.
def learn_how_to_write_better(a_poem):
   # for now, let's just ignore all other poetry
   pass
`
      }
      </Highlight>
      again, there is nothing special about the names we are giving these
      functions, they just seem to make sense given what each function does.

          </div>
          <div style={subheaderStyle}>--write</div>
          <div className={'tutorial-content'} style={contentStyle}>
          </div>
          <div style={subheaderStyle}>--critique</div>
          <div style={subheaderStyle}>--study</div>
          <div style={subheaderStyle}>using a parameters file (optional)</div>
          <div style={subheaderStyle}>putting it all together</div>
          <div style={subheaderStyle}>how to submit</div>
          <div style={subheaderStyle}>user provided tutorials</div>
          <div id='resources' style={subheaderStyle}>some resources</div>
          <div className={'tutorial-content'} style={contentStyle}>
            blah
          </div>
        </div>
      </div>
    </div>
  )
}

  // <div style={comicStyle}>
  // <Mainframe size={'medium'}/>
  // <Chip/>
  // <Floppy/>
  // <Phone/>
  // </div>
