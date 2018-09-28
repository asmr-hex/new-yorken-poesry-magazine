import argparse
import json
import random

parser = argparse.ArgumentParser()
group = parser.add_mutually_exclusive_group()
group.add_argument("--write", action="store_true")
group.add_argument("--critique", type=str, help="rate a poem between 0-1")
group.add_argument("--study", type=str, help="learn from new poems")
args = parser.parse_args()


dictionary = {
    'noun': ['people', 'food', 'data', 'theory', 'software', 'science'],
    'adj': ['colossal', 'poised', 'soupy', 'voluminous', 'zealous'],
    'verb': ['is', 'becomes', 'interprets', 'buffalos', 'overshadows'],
    'article': ['a', 'the']
}

grammar = {
    'sentence': [
        'noun verb adj noun',
        'article noun verb noun',
        'adj noun verb adj noun',
    ],
    'title': [
        'article noun',
        'adj noun',
    ],
}


def choose_grammar(grammar_type):
    return random.choice(grammar[grammar_type])


def fill_structure(structure):
    result = ''
    for pos in structure.split(' '):
        result += random.choice(dictionary[pos]) + ' '

    return result


def write_poem():
    poem = {
        'title': fill_structure(choose_grammar('title')),
        'content': fill_structure(choose_grammar('sentence')),
    }

    s = json.dumps(poem)

    print(s)


def critique_poem(poem):
    critique = {
        'score': 0.56,
    }

    s = json.dumps(critique)

    print(s)


def study_poem():
    update = {
        'success': True,
    }

    s = json.dumps(update)

    print(s)


if args.write:
    write_poem()
elif args.critique:
    critique_poem(args.critique)
elif args.study:
    study_poem()
