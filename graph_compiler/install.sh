pip install lark pre-commit

COMPILER_DIR=$(pwd)
ROOT_DIR=$(dirname $(pwd))
export PYTHONPATH=$PYTHONPATH:$COMPILER_DIR:$ROOT_DIR