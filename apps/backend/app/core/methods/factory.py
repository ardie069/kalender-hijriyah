from .umm_al_qura import UmmAlQuraMethod
from .local_hisab import LocalHisabMethod
from .local_rukyat import LocalRukyatMethod
from .ughc_base import BaseUGHCMethod

_instances = {}


def get_method_instance(method: str):
    classes = {
        "umm_al_qura": UmmAlQuraMethod,
        "local_hisab": LocalHisabMethod,
        "local_rukyat": LocalRukyatMethod,
        "ughc": BaseUGHCMethod,
    }

    if method not in classes:
        raise ValueError(f"Metode '{method}' tidak ditemukan.")

    if method not in _instances:
        _instances[method] = classes[method]()

    return _instances[method]
