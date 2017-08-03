#!/usr/bin/env python2.7

import json
from flask import (Flask,
                   make_response,
                   request)

app = Flask(__name__)

@app.route("/")
def hello():
    return "Hello World!"

@app.route("/trpe", methods=["GET", "POST"])
def trpe():
    if request.method == "POST":
        form = request.form.to_dict()
        incoming_message = form.get('incoming_message')
        command_prefix = form.get('command_prefix')
        channel = form.get('channel')
        message = trpe_demo_command(incoming_message, command_prefix, channel)
        reply = {"message": message, "status": "ok"}
        response = make_response(json.dumps(reply))
        response.headers['Content-Type'] = 'application/json'
    else:
        response = "API documentation goes here"
    return response

def trpe_demo_command(incoming_message, command_prefix, channel):
    command = incoming_message[1:]
    return 'Got message `{0!s}` with prefix `{1!s}` on channel `{2!s}`'.format(command, command_prefix, channel)

if __name__ == '__main__':
    app.run(debug=True)
