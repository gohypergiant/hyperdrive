import json


class Base:
    createdBy = None
    updatedBy = None
    deletedAt = None
    createdAt = None
    updatedAt = None

    @classmethod
    def create_from_dict(cls, dictionary):
        obj = cls()
        for key, value in dictionary.items():
            if key in cls.__dict__.keys():
                if "Hyperdrive" in str(type(value)):
                    setattr(obj, str(key), value)
                elif (type(value) is int) or (value is None):
                    setattr(obj, str(key), value)
                else:
                    try:
                        setattr(obj, str(key), json.loads(value))
                    except (json.decoder.JSONDecodeError, TypeError):
                        setattr(obj, str(key), value)
        return obj

    def extend_from_dict(self, dictionary):
        for key, value in dictionary.items():
            if key in type(self).__dict__.keys():
                setattr(self, key, value)
