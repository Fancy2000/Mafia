import helper
from flask import Flask, abort, request

app = Flask(__name__)


# Получаем отдельную запись по имени user_name
@app.route('/mafia/api/v1.0/todos/<string:user_name>', methods=['GET'])
def get_todo(user_name):
    # Получаем запись из базы данных
    response = helper.get_user(user_name)

    # Если не найдено - ошибка 404
    if response is None:
        abort(404)

    return response


# Получаем список всех пользователей коллекции
@app.route('/mafia/api/v1.0/todos', methods=['GET'])
def get_all_todos():
    return helper.get_all_users()


# Добавить элемент в коллекцию. В теле запроса должен быть передан JSON с  полями 'avatar, sex, email'
@app.route('/mafia/api/v1.0/todos', methods=['POST'])
def add_todo():
    # Если в параметрах запроса нет тела, либо нет поля 'description' - отбой
    if not request.json or not 'user_name' in request.json:
        abort(400)
    if not 'avatar' in request.json:
        abort(400)
    if not 'sex' in request.json:
        abort(400)

    if not 'email' in request.json:
        abort(400)

    # Получаем поле из запроса
    user_name = request.get_json()['user_name']
    avatar = request.get_json()['avatar']
    sex = request.get_json()['sex']
    email = request.get_json()['email']

    # Добавляем элемент в базу данных
    response = helper.add_to_list(user_name, avatar, sex, email)

    # Если не удачно - возвращаем ошибку 400
    if response is None:
        abort(400)

    # Возвращаем полное описание добавленного элемента
    return response


# Добавить пользователя в коллекцию. В теле запроса должен быть передан JSON с полями 'avatar, sex, email'
@app.route('/mafia/api/v1.0/todos/<string:user_name>', methods=['PUT'])
def update_todo(user_name):
    # Получаем запись из базы данных
    response = helper.get_user(user_name)

    # Если не найдено - ошибка 404
    if response is None:
        abort(404)

    # Если в параметрах запроса нет тела, либо нет поля 'description' - отбой
    if not request.json:
        abort(400)
    if not 'avatar' in request.json:
        abort(400)
    if not 'sex' in request.json:
        abort(400)
    if not 'email' in request.json:
        abort(400)

    # Добавляем элемент в базу данных
    response = helper.update_users(user_name, request.get_json()['avatar'], request.get_json()['sex'], request.get_json()['email'])

    # Если не удачно - возвращаем ошибку 500
    if response is None:
        abort(400)

    # Возвращаем полное описание добавленного элемента
    return response


# Удалить пользователя из коллекции
@app.route('/mafia/api/v1.0/todos/<string:user_name>', methods=['DELETE'])
def delete_task(user_name):
    # Получаем запись из базы данных
    response = helper.get_user(user_name)

    # Если не найдено - ошибка 404
    if response is None:
        abort(404)

    response = helper.remove_user(user_name)
    return response
