# Copyright 2016 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Actions to manage the system TPM."""

from glazier.lib import bitlocker
from glazier.lib.actions.base import ActionError
from glazier.lib.actions.base import BaseAction
from glazier.lib.actions.base import ValidationError


class BitlockerEnable(BaseAction):

  def Run(self):
    mode = str(self._args[0])
    try:
      bl = bitlocker.Bitlocker(mode)
      bl.Enable()
    except bitlocker.Error as e:
      raise ActionError('Failure enabling Bitlocker. (%s)' % str(e)) from e

  def Validate(self):
    self._ListOfStringsValidator(self._args, 1)
    if self._args[0] not in bitlocker.SUPPORTED_MODES:
      raise ValidationError(
          'Unknown mode for BitlockerEnable: %s' % self._args[0])

