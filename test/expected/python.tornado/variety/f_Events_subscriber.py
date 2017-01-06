#
# Autogenerated by Frugal Compiler (2.0.0-RC7)
#
# DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
#



import sys
import traceback

from thrift.Thrift import TApplicationException
from thrift.Thrift import TMessageType
from thrift.Thrift import TType
from tornado import gen
from frugal.middleware import Method
from frugal.subscription import FSubscription
from frugal.transport import TMemoryOutputBuffer

from .ttypes import *




class EventsSubscriber(object):
    """
    This docstring gets added to the generated code because it has
    the @ sign. Prefix specifies topic prefix tokens, which can be static or
    variable.
    """

    _DELIMITER = '.'

    def __init__(self, provider, middleware=None):
        """
        Create a new EventsSubscriber.

        Args:
            provider: FScopeProvider
            middleware: ServiceMiddleware or list of ServiceMiddleware
        """

        middleware = middleware or []
        if middleware and not isinstance(middleware, list):
            middleware = [middleware]
        middleware += provider.get_middleware()
        self._middleware = middleware
        self._provider = provider

    @gen.coroutine
    def subscribe_EventCreated(self, user, EventCreated_handler):
        """
        This is a docstring.
        
        Args:
            user: string
            EventCreated_handler: function which takes FContext and Event
        """

        op = 'EventCreated'
        prefix = 'foo.{}.'.format(user)
        topic = '{}Events{}{}'.format(prefix, self._DELIMITER, op)

        transport, protocol_factory = self._provider.new_subscriber()
        yield transport.subscribe(topic, self._recv_EventCreated(protocol_factory, op, EventCreated_handler))

    def _recv_EventCreated(self, protocol_factory, op, handler):
        method = Method(handler, self._middleware)

        def callback(transport):
            iprot = protocol_factory.get_protocol(transport)
            ctx = iprot.read_request_headers()
            mname, _, _ = iprot.readMessageBegin()
            if mname != op:
                iprot.skip(TType.STRUCT)
                iprot.readMessageEnd()
                raise TApplicationException(TApplicationException.UNKNOWN_METHOD)
            req = Event()
            req.read(iprot)
            iprot.readMessageEnd()
            try:
                method([ctx, req])
            except:
                traceback.print_exc()
                sys.exit(1)

        return callback




