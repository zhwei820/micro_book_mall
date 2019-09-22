# coding: utf-8
import datetime
from fabric.api import (
    local,
    cd,
    run,
    execute,
    env,
    lcd,
    task
)

env.user = 'opt'

TEST_ENV = '10.9.255.115'
PRO_ENV_NODE1 = '129.204.161.53'
PRO_ENV_NODE2 = '129.204.161.6'

ENVS = {
    'test': [{'ip': TEST_ENV, 'user': 'root'}],
    'pro': [{'ip': PRO_ENV_NODE1, 'user': 'opt'},
            {'ip': PRO_ENV_NODE2, 'user': 'opt'}],
}

SRC_DIR = '/data/opt/user-web'
TAR_DIR = '/data/opt/'


def get_passwd(env=None):
    '''
    获取开发,测试环境mysql密码
    :param env:
    :return:
    '''
    if env == TEST_ENV:
        return 'qixinchatestroot'
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
        with cd('/data/deploy_app/user-web'):
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
def create_db(db='sirius'):
    '''
    create db if db not exists
    '''
    local(
        'docker exec -it {db_env} mysql -uroot -proot -e '
        '"CREATE DATABASE IF NOT EXISTS {db} default character set utf8mb4 collate utf8mb4_unicode_ci;;"'.format(db=db, db_env=config.db_env)
    )


@task
def drop_db(db='sirius'):
    '''
    drop db
    '''
    local(
        'docker exec -it {db_env} mysql -uroot -proot -e '
        '"DROP DATABASE {db};"'.format(db=db, db_env=config.db_env)
    )


def service_restart(env='prod'):
    with cd(SRC_DIR):
        run('docker-compose pull %s' % env)


@task
def deploy_test_docker(branch='dev'):
    '''
    从gitlab 发布测试环境
    '''
    env.user = 'root'
    sync_src_code('test')
    execute(copy_exec, hosts=[TEST_ENV])
    execute(service_restart, 'test', hosts=[TEST_ENV])


@task
def deploy_pro_docker(branch='master'):
    '''
    从gitlab 发布生产环境
    '''
    sync_src_code('pro')
    execute(copy_exec, hosts=[PRO_ENV_NODE1])
    execute(service_restart, 'prod', hosts=[PRO_ENV_NODE1])
    execute(copy_exec, hosts=[PRO_ENV_NODE2])
    execute(service_restart, 'prod', hosts=[PRO_ENV_NODE2])


@task
def deploy_test():
    """
    :return:
    """
    env.user = 'root'
    local('rsync -r user-web %s:/data/opt/user-web' % (env.user + '@' + TEST_ENV))
    execute(service_restart, 'test', hosts=[TEST_ENV])
