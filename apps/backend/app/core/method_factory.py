from .methods.umm_al_qura import UmmAlQuraMethod
from .methods.local_hisab import LocalHisabMethod
from .methods.local_rukyat import LocalRukyatMethod
from .methods.ughc import UGHCMethod


def get_method_instance(method: str):
    mapping = {
        "umm_al_qura": UmmAlQuraMethod(),
        "local_hisab": LocalHisabMethod(),
        "local_rukyat": LocalRukyatMethod(),
        "ughc": UGHCMethod(),
    }

    return mapping[method]
