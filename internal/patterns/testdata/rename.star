def rename(path):
    if "skip" in path:
        return None

    parts = path.split("/")
    if len(parts) < 4:
        return path

    parts[-2] = "day=%s" % parts[-2]
    parts[-3] = "month=%s" % parts[-3]
    parts[-4] = "year=%s" % parts[-4]
    return "/".join(parts)

set_rename_filter(rename)
