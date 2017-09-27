#!/usr/bin/env python

import os

from udger import Udger
import sys

def get_ua(client_ua):
    """
    :return: dict {ua_family_code, ua_version, ua_class_code, device_class_code,
                   os_family_code, os_code}
    """
    data_dir = os.path.abspath(os.path.join(os.path.dirname( __file__ ), '..', 'data'))

    udger = Udger(data_dir)

    result = {}

    ua_obj = udger.parse_ua(client_ua)

    result['ua_family_code'] = ua_obj['ua_family_code']
    result['ua_version'] = ua_obj['ua_version']
    result['ua_class_code'] = ua_obj['ua_class_code']
    result['device_class_code'] = ua_obj['device_class_code']
    result['os_family_code'] = ua_obj['os_family_code']
    result['os_code'] = ua_obj['os_code']

    return result

if __name__ == '__main__':
    print get_ua(sys.argv[0])