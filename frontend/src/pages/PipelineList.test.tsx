/*
 * Copyright 2018 The Kubeflow Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ReactWrapper, shallow, ShallowWrapper } from 'enzyme';
import { range } from 'lodash';
import * as React from 'react';
import { RoutePage, RouteParams } from 'src/components/Router';
import { Apis } from 'src/lib/Apis';
import { ButtonKeys } from 'src/lib/Buttons';
import TestUtils from 'src/TestUtils';
import { PageProps } from './Page';
import PipelineList from './PipelineList';

describe('PipelineList', () => {
  let tree: ReactWrapper | ShallowWrapper;

  let updateBannerSpy: jest.Mock<{}>;
  let updateDialogSpy: jest.Mock<{}>;
  let updateSnackbarSpy: jest.Mock<{}>;
  let updateToolbarSpy: jest.Mock<{}>;
  let listPipelinesSpy: jest.SpyInstance<{}>;
  let listPipelineVersionsSpy: jest.SpyInstance<{}>;
  let deletePipelineSpy: jest.SpyInstance<{}>;
  let deletePipelineVersionSpy: jest.SpyInstance<{}>;

  function spyInit() {
    updateBannerSpy = jest.fn();
    updateDialogSpy = jest.fn();
    updateSnackbarSpy = jest.fn();
    updateToolbarSpy = jest.fn();
    listPipelinesSpy = jest.spyOn(Apis.pipelineServiceApiV2, 'listPipelines');
    listPipelineVersionsSpy = jest.spyOn(Apis.pipelineServiceApiV2, 'listPipelineVersions');
    deletePipelineSpy = jest.spyOn(Apis.pipelineServiceApiV2, 'deletePipeline');
    deletePipelineVersionSpy = jest.spyOn(Apis.pipelineServiceApiV2, 'deletePipelineVersion');
  }

  function generateProps(): PageProps {
    return TestUtils.generatePageProps(
      PipelineList,
      '' as any,
      '' as any,
      null,
      updateBannerSpy,
      updateDialogSpy,
      updateToolbarSpy,
      updateSnackbarSpy,
    );
  }

  async function mountWithNPipelines(n: number): Promise<ReactWrapper> {
    listPipelinesSpy.mockImplementation(() => ({
      pipelines: range(n).map(i => ({
        pipeline_id: 'test-pipeline-id' + i,
        display_name: 'test pipeline name' + i,
        name: 'test-pipeline-name' + i,
      })),
    }));
    tree = TestUtils.mountWithRouter(<PipelineList {...generateProps()} namespace='test-ns' />);
    await listPipelinesSpy;
    await TestUtils.flushPromises();
    tree.update(); // Make sure the tree is updated before returning it
    return tree;
  }

  beforeEach(() => {
    spyInit();
    jest.clearAllMocks();
  });

  afterEach(async () => {
    // unmount() should be called before resetAllMocks() in case any part of the unmount life cycle
    // depends on mocks/spies
    await tree.unmount();
    jest.restoreAllMocks();
  });

  it('renders an empty list with empty state message', () => {
    tree = shallow(<PipelineList {...generateProps()} />);
    expect(tree).toMatchSnapshot();
  });

  it('renders a list of one pipeline', async () => {
    tree = shallow(<PipelineList {...generateProps()} />);
    tree.setState({
      displayPipelines: [
        {
          created_at: new Date(2018, 8, 22, 11, 5, 48),
          description: 'test pipeline description',
          display_name: 'pipeline1',
          name: 'pipeline1',
          parameters: [],
        },
      ],
    });
    await listPipelinesSpy;
    expect(tree).toMatchSnapshot();
  });

  it('renders a list of one pipeline with no description or created date', async () => {
    tree = shallow(<PipelineList {...generateProps()} />);
    tree.setState({
      displayPipelines: [
        {
          display_name: 'pipeline1',
          name: 'pipeline1',
          parameters: [],
        },
      ],
    });
    await listPipelinesSpy;
    expect(tree).toMatchSnapshot();
  });

  it('renders a list of one pipeline with a display name that is not the same as the name', async () => {
    tree = shallow(<PipelineList {...generateProps()} />);
    tree.setState({
      displayPipelines: [
        {
          display_name: 'Pipeline One',
          name: 'pipeline1',
          parameters: [],
        },
      ],
    });
    await listPipelinesSpy;
    expect(tree).toMatchSnapshot();
  });

  it('renders a list of one pipeline with error', async () => {
    tree = shallow(<PipelineList {...generateProps()} />);
    tree.setState({
      displayPipelines: [
        {
          created_at: new Date(2018, 8, 22, 11, 5, 48),
          description: 'test pipeline description',
          error: 'oops! could not load pipeline',
          display_name: 'pipeline1',
          name: 'pipeline1',
          parameters: [],
        },
      ],
    });
    await listPipelinesSpy;
    expect(tree).toMatchSnapshot();
  });

  it('calls Apis to list pipelines, sorted by creation time in descending order', async () => {
    listPipelinesSpy.mockImplementationOnce(() => ({
      pipelines: [{ display_name: 'pipeline1', name: 'pipeline1' }],
    }));
    tree = TestUtils.mountWithRouter(<PipelineList {...generateProps()} namespace='test-ns' />);
    await listPipelinesSpy;
    expect(listPipelinesSpy).toHaveBeenLastCalledWith('test-ns', '', 10, 'created_at desc', '');
    expect(tree.state()).toHaveProperty('displayPipelines', [
      { expandState: 0, display_name: 'pipeline1', name: 'pipeline1' },
    ]);
  });

  it('has a Refresh button, clicking it refreshes the pipeline list', async () => {
    tree = await mountWithNPipelines(1);
    const instance = tree.instance() as PipelineList;
    expect(listPipelinesSpy.mock.calls.length).toBe(1);
    const refreshBtn = instance.getInitialToolbarState().actions[ButtonKeys.REFRESH];
    expect(refreshBtn).toBeDefined();
    await refreshBtn!.action();
    expect(listPipelinesSpy.mock.calls.length).toBe(2);
    expect(listPipelinesSpy).toHaveBeenLastCalledWith('test-ns', '', 10, 'created_at desc', '');
    expect(updateBannerSpy).toHaveBeenLastCalledWith({});
  });

  it('shows error banner when listing pipelines fails', async () => {
    TestUtils.makeErrorResponseOnce(listPipelinesSpy, 'bad stuff happened');
    tree = TestUtils.mountWithRouter(<PipelineList {...generateProps()} />);
    await listPipelinesSpy;
    await TestUtils.flushPromises();
    expect(updateBannerSpy).toHaveBeenLastCalledWith(
      expect.objectContaining({
        additionalInfo: 'bad stuff happened',
        message: 'Error: failed to retrieve list of pipelines. Click Details for more information.',
        mode: 'error',
      }),
    );
  });

  it('shows error banner when listing pipelines fails after refresh', async () => {
    tree = TestUtils.mountWithRouter(<PipelineList {...generateProps()} />);
    const instance = tree.instance() as PipelineList;
    const refreshBtn = instance.getInitialToolbarState().actions[ButtonKeys.REFRESH];
    expect(refreshBtn).toBeDefined();
    TestUtils.makeErrorResponseOnce(listPipelinesSpy, 'bad stuff happened');
    await refreshBtn!.action();
    expect(listPipelinesSpy.mock.calls.length).toBe(2);
    expect(listPipelinesSpy).toHaveBeenLastCalledWith(undefined, '', 10, 'created_at desc', '');
    expect(updateBannerSpy).toHaveBeenLastCalledWith(
      expect.objectContaining({
        additionalInfo: 'bad stuff happened',
        message: 'Error: failed to retrieve list of pipelines. Click Details for more information.',
        mode: 'error',
      }),
    );
  });

  it('hides error banner when listing pipelines fails then succeeds', async () => {
    TestUtils.makeErrorResponseOnce(listPipelinesSpy, 'bad stuff happened');
    tree = TestUtils.mountWithRouter(<PipelineList {...generateProps()} />);
    const instance = tree.instance() as PipelineList;
    await listPipelinesSpy;
    await TestUtils.flushPromises();
    expect(updateBannerSpy).toHaveBeenLastCalledWith(
      expect.objectContaining({
        additionalInfo: 'bad stuff happened',
        message: 'Error: failed to retrieve list of pipelines. Click Details for more information.',
        mode: 'error',
      }),
    );
    updateBannerSpy.mockReset();

    const refreshBtn = instance.getInitialToolbarState().actions[ButtonKeys.REFRESH];
    listPipelinesSpy.mockImplementationOnce(() => ({
      pipelines: [{ display_name: 'pipeline1', name: 'pipeline1' }],
    }));
    await refreshBtn!.action();
    expect(listPipelinesSpy.mock.calls.length).toBe(2);
    expect(updateBannerSpy).toHaveBeenLastCalledWith({});
  });

  it('renders pipeline names as links to their details pages', async () => {
    tree = await mountWithNPipelines(1);
    const link = tree.find('a[children="test pipeline name0"]');
    expect(link).toHaveLength(1);
    expect(link.prop('href')).toBe(
      RoutePage.PIPELINE_DETAILS_NO_VERSION.replace(
        ':' + RouteParams.pipelineId + '?',
        'test-pipeline-id0',
      ),
    );
  });

  it('always has upload pipeline button enabled', async () => {
    tree = await mountWithNPipelines(1);
    const calls = updateToolbarSpy.mock.calls[0];
    expect(calls[0].actions[ButtonKeys.NEW_PIPELINE_VERSION]).not.toHaveProperty('disabled');
  });

  it('enables delete button when one pipeline is selected', async () => {
    tree = await mountWithNPipelines(1);
    tree.find('.tableRow').simulate('click');
    expect(updateToolbarSpy.mock.calls).toHaveLength(2); // Initial call, then selection update
    const calls = updateToolbarSpy.mock.calls[1];
    expect(calls[0].actions[ButtonKeys.DELETE_RUN]).toHaveProperty('disabled', false);
  });

  it('enables delete button when two pipelines are selected', async () => {
    tree = await mountWithNPipelines(2);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(1)
      .simulate('click');
    expect(updateToolbarSpy.mock.calls).toHaveLength(3); // Initial call, then selection updates
    const calls = updateToolbarSpy.mock.calls[2];
    expect(calls[0].actions[ButtonKeys.DELETE_RUN]).toHaveProperty('disabled', false);
  });

  it('re-disables delete button pipelines are unselected', async () => {
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    expect(updateToolbarSpy.mock.calls).toHaveLength(3); // Initial call, then selection updates
    const calls = updateToolbarSpy.mock.calls[2];
    expect(calls[0].actions[ButtonKeys.DELETE_RUN]).toHaveProperty('disabled', true);
  });

  it('shows delete dialog when delete button is clicked', async () => {
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    expect(call).toHaveProperty('title', 'Delete 1 pipeline?');
  });

  it('shows delete dialog when delete button is clicked, indicating several pipelines to delete', async () => {
    tree = await mountWithNPipelines(5);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(2)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(3)
      .simulate('click');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    expect(call).toHaveProperty('title', 'Delete 3 pipelines?');
  });

  it('does not call delete API for selected pipeline when delete dialog is canceled', async () => {
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const cancelBtn = call.buttons.find((b: any) => b.text === 'Cancel');
    await cancelBtn.onClick();
    expect(deletePipelineSpy).not.toHaveBeenCalled();
  });

  it('calls delete API for selected pipeline after delete dialog is confirmed', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    expect(deletePipelineSpy).toHaveBeenLastCalledWith('test-pipeline-id0', false);
  });

  it('updates the selected indices after a pipeline is deleted', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(5);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    expect(tree.state()).toHaveProperty('selectedIds', ['test-pipeline-id0']);
    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    expect(tree.state()).toHaveProperty('selectedIds', []);
  });

  it('updates the selected indices after multiple pipelines are deleted', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(5);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(3)
      .simulate('click');
    expect(tree.state()).toHaveProperty('selectedIds', ['test-pipeline-id0', 'test-pipeline-id3']);
    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    expect(tree.state()).toHaveProperty('selectedIds', []);
  });

  it('calls delete API for all selected pipelines after delete dialog is confirmed', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(5);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(1)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(4)
      .simulate('click');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    expect(deletePipelineSpy).toHaveBeenCalledTimes(3);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id0', false);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id1', false);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id4', false);
  });

  it('shows snackbar confirmation after pipeline is deleted', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    expect(updateSnackbarSpy).toHaveBeenLastCalledWith({
      message: 'Deletion succeeded for 1 pipeline',
      open: true,
    });
  });

  it('shows error dialog when pipeline deletion fails', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    TestUtils.makeErrorResponseOnce(deletePipelineSpy, 'woops, failed');
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    const lastCall = updateDialogSpy.mock.calls[1][0];
    expect(lastCall).toMatchObject({
      content: 'Failed to delete pipeline: test-pipeline-id0 with error: "woops, failed"',
      title: 'Failed to delete some pipelines and/or some pipeline versions',
    });
  });

  it('shows error dialog when multiple pipeline deletions fail', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({ pipeline_versions: [] }));
    tree = await mountWithNPipelines(5);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(2)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(1)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(3)
      .simulate('click');
    deletePipelineSpy.mockImplementation(id => {
      if (id.indexOf(3) === -1 && id.indexOf(2) === -1) {
        // eslint-disable-next-line no-throw-literal
        throw {
          text: () => Promise.resolve('woops, failed!'),
        };
      }
    });
    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();
    // Should show only one error dialog for both pipelines (plus once for confirmation)
    expect(updateDialogSpy).toHaveBeenCalledTimes(2);
    const lastCall = updateDialogSpy.mock.calls[1][0];
    expect(lastCall).toMatchObject({
      content:
        'Failed to delete pipeline: test-pipeline-id0 with error: "woops, failed!"\n\n' +
        'Failed to delete pipeline: test-pipeline-id1 with error: "woops, failed!"',
      title: 'Failed to delete some pipelines and/or some pipeline versions',
    });

    // Should show snackbar for the one successful deletion
    expect(updateSnackbarSpy).toHaveBeenLastCalledWith({
      message: 'Deletion succeeded for 2 pipelines',
      open: true,
    });
  });

  it("delete a pipeline and some other pipeline's version together", async () => {
    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    deletePipelineVersionSpy.mockImplementation(() => Promise.resolve());

    // Mock listPipelineVersions with different responses based on pipeline ID
    listPipelineVersionsSpy.mockImplementation((pipelineId: string) => {
      if (pipelineId === 'test-pipeline-id1') {
        // Return versions for expansion functionality
        return {
          pipeline_versions: [
            {
              display_name: 'test-pipeline-id1_name',
              name: 'test-pipeline-id1_name',
              pipeline_id: 'test-pipeline-id1',
              pipeline_version_id: 'test-pipeline-version-id1',
            },
          ],
        };
      }
      // Return empty for other pipelines so dialog shows "Delete" instead of "Delete All"
      return { pipeline_versions: [] };
    });

    tree = await mountWithNPipelines(2);
    tree
      .find('button[aria-label="Expand"]')
      .at(1)
      .simulate('click');
    await listPipelineVersionsSpy;
    tree.update();

    // select pipeline of id 'test-pipeline-id0'
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    // select pipeline version of id 'test-pipeline-id1_default_version' under pipeline 'test-pipeline-id1'
    tree
      .find('.tableRow')
      .at(2)
      .simulate('click');

    expect(tree.state()).toHaveProperty('selectedIds', ['test-pipeline-id0']);
    expect(tree.state()).toHaveProperty('selectedVersionIds', {
      'test-pipeline-id1': ['test-pipeline-version-id1'],
    });

    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    const call = updateDialogSpy.mock.calls[0][0];
    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete');
    await confirmBtn.onClick();

    await deletePipelineSpy;
    await deletePipelineVersionSpy;

    expect(deletePipelineSpy).toHaveBeenCalledTimes(1);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id0', false);

    expect(deletePipelineVersionSpy).toHaveBeenCalledTimes(1);
    expect(deletePipelineVersionSpy).toHaveBeenCalledWith(
      'test-pipeline-id1',
      'test-pipeline-version-id1',
    );

    expect(tree.state()).toHaveProperty('selectedIds', []);
    expect(tree.state()).toHaveProperty('selectedVersionIds', { 'test-pipeline-id1': [] });

    // Should show snackbar for the one successful deletion
    expect(updateSnackbarSpy).toHaveBeenLastCalledWith({
      message: 'Deletion succeeded for 1 pipeline and 1 pipeline version',
      open: true,
    });
  });

  it('shows "Delete All" dialog when pipelines have versions and calls API with cascade=true', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({
      pipeline_versions: [
        {
          pipeline_version_id: 'test-version-id-1',
          display_name: 'test version 1',
          name: 'test-version-1',
        },
        {
          pipeline_version_id: 'test-version-id-2',
          display_name: 'test version 2',
          name: 'test-version-2',
        },
      ],
    }));

    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');

    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    await TestUtils.flushPromises();

    const call = updateDialogSpy.mock.calls[0][0];

    // Verify the dialog shows cascade warning and "Delete All" button
    expect(call.title).toBe('Delete 1 pipeline?');
    expect(call.content).toContain('pipeline has existing versions');
    expect(call.content).toContain('Deleting this pipeline will also delete all its versions');
    expect(call.content).toContain('This action cannot be undone');

    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete All');
    expect(confirmBtn).toBeDefined();

    await confirmBtn.onClick();
    expect(deletePipelineSpy).toHaveBeenLastCalledWith('test-pipeline-id0', true);
    expect(updateSnackbarSpy).toHaveBeenLastCalledWith({
      message: 'Deletion succeeded for 1 pipeline',
      open: true,
    });
  });

  it('shows "Delete All" dialog for multiple pipelines with versions', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({
      pipeline_versions: [
        {
          pipeline_version_id: 'test-version-id-1',
          display_name: 'test version 1',
          name: 'test-version-1',
        },
      ],
    }));

    deletePipelineSpy.mockImplementation(() => Promise.resolve());
    tree = await mountWithNPipelines(3);

    // Select multiple pipelines
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');
    tree
      .find('.tableRow')
      .at(1)
      .simulate('click');

    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    await TestUtils.flushPromises();

    const call = updateDialogSpy.mock.calls[0][0];

    // Verify the dialog shows cascade warning for multiple pipelines
    expect(call.title).toBe('Delete 2 pipelines?');
    expect(call.content).toContain('pipelines have existing versions');
    expect(call.content).toContain('Deleting these pipelines will also delete all their versions');

    const confirmBtn = call.buttons.find((b: any) => b.text === 'Delete All');
    expect(confirmBtn).toBeDefined();
    await confirmBtn.onClick();
    expect(deletePipelineSpy).toHaveBeenCalledTimes(2);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id0', true);
    expect(deletePipelineSpy).toHaveBeenCalledWith('test-pipeline-id1', true);
  });

  it('does not call delete API when "Delete All" dialog is canceled', async () => {
    listPipelineVersionsSpy.mockImplementation(() => ({
      pipeline_versions: [
        {
          pipeline_version_id: 'test-version-id-1',
          display_name: 'test version 1',
          name: 'test-version-1',
        },
      ],
    }));

    tree = await mountWithNPipelines(1);
    tree
      .find('.tableRow')
      .at(0)
      .simulate('click');

    const deleteBtn = (tree.instance() as PipelineList).getInitialToolbarState().actions[
      ButtonKeys.DELETE_RUN
    ];
    await deleteBtn!.action();
    await TestUtils.flushPromises();

    const call = updateDialogSpy.mock.calls[0][0];
    const cancelBtn = call.buttons.find((b: any) => b.text === 'Cancel');
    await cancelBtn.onClick();
    expect(deletePipelineSpy).not.toHaveBeenCalled();
  });
});
