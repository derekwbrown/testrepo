from collections import defaultdict
import sys
import time
import win32pdh
import _winreg


pdh_counter_dict = defaultdict(list)
pdh_help_dict = {}

def _get_counter_dictionaries():
    if pdh_counter_dict:
        # already populated
        return

    try:
        val, t = _winreg.QueryValueEx(_winreg.HKEY_PERFORMANCE_DATA, "Counter 009")
    except:
        raise

    # val is an array of strings.  The underlying win32 API returns a list of strings
    # which is the counter name, counter index, counter name, counter index (in windows,
    # a multi-string value)
    #
    # the python implementation translates the multi-string value into an array of strings.
    # the array of strings then becomes
    # array[0] = counter_index_1
    # array[1] = counter_name_1
    # array[2] = counter_index_2
    # array[3] = counter_name_2
    #
    # see https://support.microsoft.com/en-us/help/287159/using-pdh-apis-correctly-in-a-localized-language
    # for more detail

    # create a table of the keys to the counter index, because we want to look up
    # by counter name.
    idx = 0
    idx_max = len(val)
    while idx < idx_max:
        # counter index is idx , counter name is idx + 1
        pdh_counter_dict[val[idx+1]].append(val[idx])
        idx += 2

    k = _winreg.OpenKey(_winreg.HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Perflib\\CurrentLanguage")
    val, t = _winreg.QueryValueEx(k, "Help")
    idx = 0
    idx_max = len(val)
    while idx < idx_max:
        pdh_help_dict[int(val[idx])] = val[idx + 1]
        idx += 2


def get_help_for_metric(metric):
    res = []
    indexes = pdh_counter_dict[metric]
    #print "%s indexes %s" % (metric, str(indexes))
    for i in indexes:
        ii = int(i) +1
        if ii in pdh_help_dict:
            res.append(pdh_help_dict[ii])
    return res


def enumerate_counter_class(class_name):
    class_name_index_list = pdh_counter_dict[class_name]
    if len(class_name_index_list) == 0:
        print("Class %s was not in counter name list, attempting english counter" % class_name)
        _class_name = class_name
    else:
        if len(class_name_index_list) > 1:
            print("Class %s had multiple (%d) indices, using first" % (class_name, len(class_name_index_list)))
        _class_name = win32pdh.LookupPerfNameByIndex(None, int(class_name_index_list[0]))

    counters, instances = win32pdh.EnumObjectItems(None, None, _class_name, win32pdh.PERF_DETAIL_WIZARD)
    for c in counters:
        helps = get_help_for_metric(c)
        for h in helps:
            print(" : %s\%s : %s" % (_class_name, c, h))

    
_get_counter_dictionaries()
#for k in pdh_help_dict:
#    print("key %d string %s\n" % (k, pdh_help_dict[k]))


#metrics = sys.argv[1:]
cl = sys.argv[1]
enumerate_counter_class(cl)
#print "cl %s\n" % cl
#for h in pdh_counter_dict[cl]:
#    print("%s\n" % h)
#    help = get_help_for_metric(h)
#    print("Help: %s\n" %help)
#for metric in metrics:
