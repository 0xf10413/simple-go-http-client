#!/usr/bin/env python3
# Copyright flo <flo@knightknight>
# Distributed under terms of the MIT license.

from flask import Flask, escape, request, session, redirect, url_for
from json import dumps
from collections import namedtuple

Proxy = namedtuple("Proxy", "Name IP Port")
Database = namedtuple("Database", "Name IP Port")

app = Flask(__name__)
app.secret_key = 'blahblahblahsecret'

@app.route('/')
def index():
    if 'username' in session:
        return "Logged in as %s" % escape(session['username'])
    return "Not logged in"

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        session['username'] = request.form['username']
        return redirect(url_for('index'))
    return '''
            <form method="post">
              <p><input type=text name=username>
              <p><input type=submit value=Login>
            </form>'''

@app.route('/get-proxies')
def get_proxies():
    if 'username' not in session:
        return redirect(url_for("login"))

    proxies = [Proxy(Name=f"Proxy-{i}", IP="localhost", Port=f"{50051+i}")._asdict()
        for i in range(6)]
    return dumps(proxies)

@app.route('/get-databases')
def get_databases():
    if 'username' not in session:
        return redirect(url_for("login"))

    if request.method == 'POST':
        session['username'] = request.form['username']
        return redirect(url_for('index'))
    databases = [Database(Name=f"Database-{i}", IP="localhost", Port=f"{50051+i}")._asdict()
        for i in range(4)]
    return dumps(databases)
