import {get, map} from 'lodash'


export const getPoemsByIssueVolume = (volume, state) =>
  map(
    get(state, `issuesByVolume.${volume}.poems`, []),
    id => {
      // join author into poem
      const poem = get(state, `poems.${id}`, {})

      return {
        ...poem,
        author: get(state, `poets.${poem.author.id}`, {})
      }
    },
    [],
  )
