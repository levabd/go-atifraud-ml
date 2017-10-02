#!/usr/bin/env python

import os

from udger import Udger
import sys
def is_crawler(client_ip, client_ua):
    """
    :return: crawler or not
    """

    bots_ua_family = {
        # Search engine or antivirus or SEO bots
        'googlebot',
        'siteexplorer',
        'sputnikbot',
        'bingbot',
        'mj12bot',
        'yandexbot',
        'cliqzbot',
        'avast_safezone',
        'megaindex',
        'genieo_web_filter',
        'uptimebot',
        'ahrefsbot',
        'wordpress_pingback',
        'admantx_platform_semantic_analyzer',
        'leikibot',
        'mnogosearch',
        'safednsbot',
        'easybib_autocite',
        'sogou_spider',
        'surveybot',
        'baiduspider',
        'indy_library',
        'mail-ru_bot',
        'pocketparser',
        'virustotal',
        'feedfetcher_google',
        'virusdie_crawler',
        'surdotlybot',
        'yoozbot',
        'facebookbot',
        'linkdexbot',
        'prlog',
        'thinglink_imagebot',
        'obot',
        'spyonweb',
        'easybib_autocite',
        'avira_crawler',
        'pulsepoint_xt3_web_scraper',
        'comodospider',
        'girafabot',
        'avira_scout',
        'salesintelligent',
        'kaspersky_bot',
        'xenu',
        'maxpointcrawler',
        'seznambot',
        'magpie-crawler',
        'yesupbot',
        'startmebot',
        'brandprotect_bot',
        'ask_jeeves-teoma',
        'duckduckgo_app',
        'linqiabot',
        'flipboardbot',
        'cat_explorador',
        'huaweisymantecspider',
        'coccocbot', 
        'powermarks', 
        'prism', 
        'leechcraft', 
        'wkhtmltopdf',

        # I think next is potencial bad bot (framework for apps or bad crowler)
        'java',
        'www::mechanize',
        'grapeshotcrawler',
        'netestate_crawler',
        'ccbot',
        'plukkie',
        'metauri',
        'silk',
        'phantomjs',
        'python-requests',
        'okhttp',
        'python-urllib',
        'netcraft_crawler',
        'go_http_package',
        'google_app',
        'android_httpurlconnection',
        'curl',
        'w3m',
        'wget',
        'getintentcrawler',
        'scrapy',
        'crawler4j',
        'apache-httpclient',
        'feedparser',
        'php',
        'simplepie',
        'lwp::simple',
        'libwww-perl',
        'apache_synapse',
        'scrapy_redis',
        'winhttp',
        'johnhew_crawler',
        'poe-component-client-http',
        'joc_web_spider',

        #Text Browsers
        'elinks',
        'links',
        'lynx'
    }

    data_dir = os.path.abspath(os.path.join(os.path.dirname( __file__ ), '..', '..', 'data', 'db' ))
    udger = Udger(data_dir)

    ua_c = udger.parse_ua(client_ua)
    print(ua_c)
    if ((udger.parse_ip(client_ip)['ip_classification_code'] == 'crawler') or
        (ua_c['ua_class_code'] == 'crawler') or (ua_c['ua_family_code'] in bots_ua_family)):
        return True
    else:
        #print("{} - {}".format(ua_c['ua_family_code'], ua_c['ua_family_code'] in bots_ua_family))
        return False

if __name__ == '__main__':
    print (is_crawler(sys.argv[1], sys.argv[2]))