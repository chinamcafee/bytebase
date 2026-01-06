<template>
  <!-- DISABLED: Product improvement settings disabled to prevent metrics collection -->
  <div class="py-6 lg:flex">
    <div class="text-left lg:w-1/4">
      <h1 class="text-2xl font-bold text-gray-400">
        {{ $t("settings.general.workspace.product-improvement.self") }}
      </h1>
    </div>
    <div class="flex-1 lg:px-4">
      <div class="text-sm text-gray-400">
        <p>Metrics collection has been disabled in this installation to prevent communication with external servers.</p>
        <p class="mt-2">The original setting allowed sharing anonymized usage data to help improve Bytebase.</p>
      </div>
    </div>
  </div>

  <!-- The original form is disabled:
  <div id="product-improvement" class="py-6 lg:flex">
    <div class="text-left lg:w-1/4">
      <h1 class="text-2xl font-bold">
        {{ $t("settings.general.workspace.product-improvement.self") }}
      </h1>
    </div>
    <div class="flex-1 lg:px-4">
      <div class="flex items-center gap-x-2">
        <Switch
          v-model:value="state.enableMetricCollection"
          :text="true"
          :disabled="!allowEdit"
        />
        <span class="text-sm">
          {{ $t("settings.general.workspace.product-improvement.participate") }}
        </span>
      </div>
      <div class="mt-1 text-sm text-gray-400">
        {{ $t("settings.general.workspace.product-improvement.description") }}
      </div>
    </div>
  </div>
  -->
</template>

<script lang="ts" setup>
// DISABLED: Product improvement settings disabled to prevent metrics collection
// All functionality related to metrics collection has been disabled.

import { computed } from "vue";

defineProps<{
  allowEdit: boolean;
}>();

// Since metrics collection is disabled, this component no longer needs to manage state
const isDirty = computed(() => false);
const revert = () => {
  // No state to revert since metrics collection is disabled
};
const update = async () => {
  // No updates needed since metrics collection is disabled
};

defineExpose({
  isDirty,
  revert,
  update,
  title: "Product Improvement",
});

// The original script is disabled:
/*
import { create } from "@bufbuild/protobuf";
import { FieldMaskSchema } from "@bufbuild/protobuf/wkt";
import { isEqual } from "lodash-es";
import { computed, reactive } from "vue";
import { Switch } from "@/components/v2";
import { useSettingV1Store } from "@/store";
import type { WorkspaceProfileSetting } from "@/types/proto-es/v1/setting_service_pb";

interface LocalState {
  enableMetricCollection: boolean;
}

const settingV1Store = useSettingV1Store();

const getInitialState = (): LocalState => {
  return {
    enableMetricCollection:
      settingV1Store.workspaceProfileSetting?.enableMetricCollection ?? true,
  };
};

const state = reactive(getInitialState());
const originalState = reactive(getInitialState());

const isDirty = computed(() => {
  return !isEqual(state, originalState);
});

const revert = () => {
  Object.assign(state, originalState);
};

const update = async () => {
  const payload: Partial<WorkspaceProfileSetting> = {
    enableMetricCollection: state.enableMetricCollection,
  };

  await settingV1Store.updateWorkspaceProfile({
    payload,
    updateMask: create(FieldMaskSchema, {
      paths: ["value.workspace_profile.enable_metric_collection"],
    }),
  });

  Object.assign(originalState, state);
};
*/
</script>
