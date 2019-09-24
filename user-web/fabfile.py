# coding: utf-8
import datetime

from fabric.api import (
    local,
    cd,
    run,
    execute,
    env,
    task
)

env.user = 'opt'

TEST_NODE1 = '10.9.255.115'
TEST_NODE1 = '172.17.0.8'

PRO_ENV_NODE1 = '129.204.161.53'
PRO_ENV_NODE2 = '129.204.161.6'

SRC_DIR = '/data/opt/user-web'
TAR_DIR = '/data/opt/'


def get_supervisor_env(env='test'):
    if env.lower() == 'test':
        return 'export ConfigAddress=10.35.195.56:9600;'
    if env.lower() == 'prod':
        return 'export ConfigAddress=10.35.195.56:9600;'
    return ''


def get_passwd(env='test'):
    '''
    获取开发,测试环境mysql密码
    :param env:
    :return:
    '''
    if env.lower() == 'test':
        return 'rootroot'
    return 'root'


class Config(object):
    db_env = 'product_mysql'


config = Config()


@task
def backup(envip=None):
    '''
    备份数据库, 开发\测试环境
    '''
    now = datetime.datetime.now().strftime('%Y%m%d%H%M')
    exec_str = """
            docker exec {db_env} sh -c \
            'mysqldump -uroot -p{p} user-web ' \
            > ./backup/data_{0}.sql; \
            cp ./backup/data_{0}.sql ./backup/data.sql
            """.format(now, db_env=config.db_env, p=get_passwd(envip))

    if not envip:
        local('mkdir ./backup -p')
        local(exec_str)
    else:
        with cd('/data/opt/user-web'):
            run('mkdir ./backup -p')
            run(exec_str)


@task
def recovery(envip=None, sql_file=None):
    '''
    恢复数据库
    '''
    filename = './backup/data.sql' if not sql_file else sql_file
    drop_db()
    create_db()
    local(
        'docker exec -i {db_env} mysql -uroot -p{p} '
        'user-web < {filename}'.format(filename=filename, db_env=config.db_env, p=get_passwd(envip))
    )


@task
def create_db(db='micro_book_mall'):
    '''
    create db if db not exists
    '''
    local(
        'docker exec -it {db_env} mysql -uroot -proot -e '
        '"CREATE DATABASE IF NOT EXISTS {db} default character set utf8mb4 collate utf8mb4_unicode_ci;;"'.format(db=db,
                                                                                                                 db_env=config.db_env)
    )


@task
def drop_db(db='micro_book_mall'):
    '''
    drop db
    '''
    local(
        'docker exec -it {db_env} mysql -uroot -proot -e '
        '"DROP DATABASE {db};"'.format(db=db, db_env=config.db_env)
    )


def service_restart(env='prod'):
    run('mkdir -p /data/user-web/log')
    with cd(SRC_DIR):
        run('%s supervisord -c /data/opt/user-web/deploy/supervisord.conf | true ' % get_supervisor_env(env))
        run('supervisorctl -c /data/opt/user-web/deploy/supervisord.conf restart user_web ')


def make_build():
    local('make build')


def sync_exe(node):
    local('rsync -r user-web %s:/data/opt/user-web/' % (env.user + '@' + node))
    local('rsync -r deploy %s:/data/opt/user-web/' % (env.user + '@' + node))


@task
def deploy_pro():
    '''
    '''
    env.user = 'opt'
    for node in (PRO_ENV_NODE1, PRO_ENV_NODE2):
        make_build()
        sync_exe(node)
        execute(service_restart, 'prod', hosts=[node])


@task
def deploy_test():
    '''
    '''
    env.user = 'root'
    for node in (TEST_NODE1):
        make_build()
        sync_exe(node)
        execute(service_restart, 'prod', hosts=[node])


@task
def deploy_test_from_local():
    """
    """
    env.user = 'root'
    sync_exe(TEST_NODE1)
    execute(service_restart, 'test', hosts=[TEST_NODE1])
