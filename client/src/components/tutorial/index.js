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
  const tocStyle = {
    color: '#c3a5d3',
    fontSize: '0.9em',
    fontWeight: 'normal',
    fontStyle: 'normal',
    textShadow: '3px 3px #edc9fc',
    display: 'flex',
    justifyContent: 'space-between',
    width: '60%',
  }
  const tocContainerStyle = {
    display: 'flex',
    justifyContent: 'center'
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
          <div style={subheaderStyle}>contents</div>
          <div style={{...contentStyle, ...keywordStyle, display: 'flex', flexDirection: 'column', justifyContent: 'center', textAlign: 'center'}}>
            <div style={tocContainerStyle}><a style={tocStyle} href='#introduction'><div>intr0duction</div><div>֎</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#whats-a-poet'><div>what's an algorithmic poet?</div><div>⎈</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#first-poet'><div>your first poet</div><div>♥</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#more-interesting-poet'><div>a more interesting poet</div><div>★</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#refining-poet-taste'><div>refining your poet's taste</div><div>☂</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#learning-from-new-poetry'><div>learning from new poetry</div><div>☻</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#parameters-file'><div>using the parameters file</div><div>☾</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#how-to-submit'><div>how to submit</div><div>♫</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#user-tutorials'><div>user submitted tutorials</div><div>⚑</div></a></div>
            <div style={tocContainerStyle}><a style={tocStyle} href='#resources'><div>some resources</div><div>⚛</div></a></div>
          </div>
          <div id='introduction' style={subheaderStyle}>intr0duction</div>
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
          <div id='whats-a-poet' style={subheaderStyle}>what's an algorithmic poet?</div>
          <div className={'tutorial-content'} style={contentStyle}>
            an <em>algorithmic poet</em> is a computer program you create that
            can <em style={keywordStyle}>write poems</em>, <em style={keywordStyle}>critique poems</em>,
      & <em style={keywordStyle}>learn how to write better poems</em>! your program
    should be able to handle each of the following command-line arguments in order to perform these tasks.
      <div style={{display: 'flex', justifyContent: 'center'}}>
      <ul style={{width: '70%', listStyleType: 'none'}}>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--write</span> generate and print a poem to standard out (stdout). the ouput
    must be printed in json, resembling <span style={{...codeStyle, fontSize:'0.6em', backgroundColor: '#ffffff'}}>{String.raw`{"title": "<your-title>", "content": "<your-content>"}`}</span>
      </li>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--critique POEM</span> critique a given poem, POEM, and print a score between
    0 (worst) and 1 (best). again, output should be formatted in json
    like <span style={{...codeStyle, fontSize:'0.6em', backgroundColor: '#ffffff'}}>{String.raw`{"score": <your-score>}`}</span>
      </li>
      <li style={{marginBottom: '0.7em'}}>
      <span style={codeStyle}>--study POEM</span> read a given poem, POEM, and <em>optionally</em> use
    it to modify how your program writes poetry. this task should not print any output to stdout. yet again, json
    ouput <span style={{...codeStyle, fontSize:'0.6em', backgroundColor: '#ffffff'}}>{String.raw`{"success": <true|false>}`}</span>
      </li>
      </ul>
      </div>
      for example, if you write a poet in Python called <span style={codeStyle}>poet_bot.py</span>,
    then you can ask it to perform each task by running it on the command line like this,
      <Highlight className="shell">
      {
        String.raw`$ python poet_bot.py --write  # write a poem
{"title": "flowers", "content": "roses are red, violets are blue"}

$ python poet_bot.py --critique "this is a bad poem" # critique a poem
{"score": 0.88}

$ python poet_bot.py --study "this is a poem to study"  # learn how to write better
{"success": true}
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
          <div id='first-poet' style={subheaderStyle}>your first poet</div>
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
            development environment (IDE), type out the following in your Python file:

            <Highlight className="python poet-body-code">
              {
                String.raw`# poet_bot.py
import json


# this function generates and prints a poem.
def write_poem():
    poem = {
             'title': 'flowers',
             'content': 'roses are red, violets are blue',
    }
    print(json.dumps(poem))


# when given a poem, this function critiques it on
# a scale from 0-1 and prints the score.
def critique_poem(a_poem):
   print(json.dumps({'score': 0.88}))


# when given a poem, this function can allow our
# program to potentially learn new styles or approaches to
# writing poetry.
def learn_how_to_write_better(a_poem):
   # for now, let's just ignore all other poetry
   # and say that we successfully learned even tho we didn't
   # actually do anything lol
   print(json.dumps({'success': True}))
`
              }
            </Highlight>
            again, there is nothing special about the names we are giving these
            functions other than telling us clearly what the functions do.
            <br/><br/>
            this looks pretty good, but there is something missing. we need to be able to
            call the right function when our python file is run with each command-line argument.
            so when we use the <span style={codeStyle}>--write</span> command-line argument, we
            must run our <span style={codeStyle}>write_poem()</span> function. and so on.
            <br/><br/>
            luckily there is a builtin python module called <span style={codeStyle}>argparse</span> for
            doing just this. let's import it and include some logic to <em>parse</em>, or read, the
            three required command-line arguments.
            <Highlight className="python poet-body-code">
              {
                String.raw`# poet_bot.py
import json
import argparse


# create a new argument parser
parser = argparse.ArgumentParser()

# if --write is given, this sets the args.write variable to True,
# otherwise its set to False
parser.add_argument('--write', action='store_true')

# if the --critique POEM arguments are given, this stores the 
# string POEM in args.critique, otherwise args.critique is None
parser.add_argument('--critique', type=str, help='rate a poem between 0-1')

# if the --study POEM arguments are given, this stores the 
# string POEM in args.study, otherwise args.study is None
parser.add_argument('--study', type=str, help='learn from new poems')

# get command-line arguments we described above and store
# them in the args variable.
args = parser.parse_args()


# this function generates and prints a poem.
def write_poem():
    poem = {
     'title': 'flowers',
     'content': 'roses are red, violets are blue',
    }
    print(json.dumps(poem))


# when given a poem, this function critiques it on
# a scale from 0-1 and prints the score.
def critique_poem(a_poem):
   print(json.dumps({'score': 0.88}))


# when given a poem, this function can allow our
# program to potentially learn new styles or approaches to
# writing poetry.
def learn_how_to_write_better(a_poem):
   # for now, let's just ignore all other poetry
   # and say that we successfully learned even tho we didn't
   # actually do anything lol
   print(json.dumps({'success': True}))


# run functions according to which command-line arguments were given
if args.write:
  # bingo! the --write argument was given
  write_poem()
elif args.critique:
  # wahoo! the --critique POEM arguments were given
  critique_poem(args.critique)
elif args.study:
  # awoo! the --study POEM arguments were given
  study_poem(args.study)
`
            }
    </Highlight>
      now we have a section of code at the beginning which configures <em>how we will parse</em> the
    command-line arguments given and we have a section of code at the end which uses the parsed
    arguments to <em>decide which function we will call</em>.
      <br/><br/>
    and that's it. we have our first poet!
    like we said, it isn't very interesting, but it implements the example we gave in the previous section
    and is a perfectly valid poet.
      <br/><br/>
      in the following sections, we will elaborate on this poet to make it more interesting.
      </div>
      
      <div id='more-interesting-poet'style={subheaderStyle}>a more interesting poet</div>
      <div className={'tutorial-content'} style={contentStyle}>
      wow, we've already written our first poet! why don't you take a <em>five minute</em> break
    to stretch, drink some water, say hi to your friends, water your plants.
      <br/><br/><br/>
      <div style={{textAlign: 'center'}}>. . .</div>
      <br/><br/>
      hello again, good to have you back! if you recall, our first poet is able to perform all the
    required tasks, but there is one problem: it outputs the same results everytime. now this isn't
    very interesting. but we can change this! we believe in you!
      <br/><br/>
      there are many approaches we could take, maybe an infinite number. for example,
    </div>
      
      <div id='refining-poet-taste' style={subheaderStyle}>refining your poet's taste</div>
      <div className={'tutorial-content'} style={contentStyle}>
      
    </div>
      <div id='learning-from-new-poetry' style={subheaderStyle}>learning from new poetry</div>
      <div className={'tutorial-content'} style={contentStyle}>
      
    </div>
      <div id='parameters-file' style={subheaderStyle}>using a parameters file (optional)</div>
      <div className={'tutorial-content'} style={contentStyle}>
      
    </div>
      <div id='how-to-submit' style={subheaderStyle}>how to submit</div>
      <div className={'tutorial-content'} style={contentStyle}>
      
    </div>
      <div id='user-tutorials' style={subheaderStyle}>user provided tutorials</div>
      <div className={'tutorial-content'} style={contentStyle}>
      
    </div>
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
