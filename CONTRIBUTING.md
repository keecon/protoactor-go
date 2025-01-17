# Contributing to Proto.Actor

In case of questions about the contribution process or for discussion of specific issues please visit
the [Proto.Actor gitter chat](https://gitter.im/asynkron/gam).

# Navigating around the project & codebase

## Branches summary

Depending on which version (or sometimes module) you want to work on, you should target a specific branch as explained
below:

- `dev` – active development branch of Proto.Actor beta.

## Tags

Proto.Actor uses tags to categorise issues into groups or mark their phase in development.

Most notably many tags start with a `t:` prefix (as in `topic:`), which categorises issues in terms of which module they
relate to. Examples are:

- [t:actor](https://github.com/keecon/protoactor-go/labels/t%3Aactor)
- [t:remote](https://github.com/keecon/protoactor-go/labels/t%3Aremote)
- [t:cluster](https://github.com/keecon/protoactor-go/labels/t%3Acluster)
- [t:router](https://github.com/keecon/protoactor-go/labels/t%3Arouter)
- [See all labels](https://github.com/keecon/protoactor-go/labels)

In general _all issues are open for anyone working on them_, however if you're new to the project and looking for an
issue
that will be accepted and likely is a nice one to get started you should check out the following tags:

- [community](https://github.com/keecon/protoactor-go/labels/community) - which identifies issues that the core team
  will likely not have time to work on, or the issue is a nice entry level ticket. If you're not sure how to solve a
  ticket but would like to work on it feel free to ask in the issue about clarification or tips.
- [nice-to-have (low-priority)](https://github.com/keecon/protoactor-go/labels/nice-to-have) - are tasks which make
  sense, however are not very high priority (in case of other very high priority issues). If you see something
  interesting in this list, a contribution would be really wonderful!

Another group of tickets are those which start from a number. They're used to signal in what phase of development an
issue is:

- [0 - new](https://github.com/Proto.Actor/Proto.Actor/labels/0%20-%20new) - is assigned when a ticket is unclear on its
  purpose or if it is valid or not. Sometimes the additional tag `discuss` is used to mark such tickets, if they propose
  large scale changes and need more discussion before moving into triaged (or being closed as invalid).
- [1 - triaged](https://github.com/Proto.Actor/Proto.Actor/labels/1%20-%20triaged) - roughly speaking means "this ticket
  makes sense". Triaged tickets are safe to pick up for contributing in terms of likeliness of a patch for it being
  accepted. It is not recommended to start working on a ticket that is not triaged.
- [2 - pick next](https://github.com/Proto.Actor/Proto.Actor/labels/2%20-%20pick%20next) - used to mark issues which are
  next up in the queue to be worked on. Sometimes it's also used to mark which PRs are expected to be reviewed/merged
  for the next release. The tag is non-binding, and mostly used as an organisational helper.
- [3 - in progress](https://github.com/Proto.Actor/Proto.Actor/labels/3%20-%20in%20progress) - means someone is working
  on this ticket. If you see a ticket that has the tag, however seems inactive, it could have been an omission with
  removing the tag, feel free to ping the ticket then if it's still being worked on.

The last group of special tags indicate specific states a ticket is in:

- [bug](https://github.com/keecon/protoactor-go/labels/bug) - bugs take priority in being fixed above features. The
  core team dedicates a number of days to working on bugs each sprint. Bugs which have reproducers are also great for
  community contributions as they're well isolated. Sometimes we're not as lucky to have reproducers though, then a
  bugfix should also include a test reproducing the original error along with the fix.

# Proto.Actor contributing guidelines

These guidelines apply to all Proto.Actor projects, by which we mean both the `keecon/protoactor-go` repository,
as well as any plugins or additional repositories located under the asynkron/protoactor\*.

These guidelines are meant to be a living document that should be changed and adapted as needed.
We encourage changes that make it easier to achieve our goals in an efficient way.

## General workflow

The steps below describe how to get a patch into a main development branch (e.g. `dev`).
The steps are exactly the same for everyone involved in the project (be it core team, or first time contributor).

1. Make sure an issue exists in the [issue tracker](https://github.com/keecon/protoactor-go/issues) for the work you
   want to contribute.
   - If there is no ticket for it, [create one](https://github.com/keecon/protoactor-go/issues/new) first.
1. [Fork the project](https://github.com/keecon/protoactor-go#fork-destination-box) on GitHub. You'll need to create a
   feature-branch for your work on your fork, as this way you'll be able to submit a pull request against the mainline
   Proto.Actor.
1. Create a branch on your fork and work on the feature. For
   example: `git checkout -b wip-custom-serialization-protoactor`
   - Please make sure to follow the general quality guidelines (specified below) when developing your patch.
   - Please write additional tests covering your feature and adjust existing ones if needed before submitting your pull
     request.
1. Once your feature is complete, prepare the commit following
   our [Creating Commits And Writing Commit Messages](#creating-commits-and-writing-commit-messages). For example, a
   good commit message would be: `Adding compression support for Manifests #22222` (note the reference to the ticket it
   aimed to resolve).
1. If it's a new feature, or a change of behaviour, document it on
   the [Proto.Actor-docs](https://github.com/asynkron/gam-web/tree/master/src/docs), remember, an undocumented feature
   is not a feature.
1. Now it's finally time to [submit the pull request](https://help.github.com/articles/using-pull-requests)!
1. If you have not already done so, you will be asked by our CLA bot
   to [sign the Asynkron CLA](https://cla-assistant.io/keecon/protoactor-go) online. CLA stands for Contributor
   License Agreement and is a way of protecting intellectual property disputes from harming the project.
1. Now both committers and interested people will review your code. This process is to ensure the code we merge is of
   the best possible quality, and that no silly mistakes slip though. You're expected to follow-up these comments by
   adding new commits to the same branch. The commit messages of those commits can be more lose, for
   example: `Removed debugging using printline`, as they all will be squashed into one commit before merging into the
   main branch.
   - The community and team are really nice people, so don't be afraid to ask follow up questions if you didn't
     understand some comment, or would like clarification on how to continue with a given feature. We're here to help,
     so feel free to ask and discuss any kind of questions you might have during review!
1. After the review you should fix the issues as needed (pushing a new commit for new review etc.), iterating until the
   reviewers give their thumbs up–which is signalled usually by a comment saying `LGTM`, which means "Looks Good To Me".
   - In general a PR is expected to get 2 LGTMs from the team before it is merged. If the PR is trivial, or under
     special circumstances (such as most of the team being on vacation, a PR was very thoroughly reviewed/tested and
     surely is correct) one LGTM may be fine as well.
1. Once everything is said and done, your pull request gets merged :tada: Your feature will be available with the next
   “earliest” release milestone (i.e. if back-ported so that it will be in release x.y.z, find the relevant milestone
   for that release). And of course you will be given credit for the fix in the release stats during the release's
   announcement. You've made it!

The TL;DR; of the above very precise workflow version is:

1. Fork Proto.Actor
2. Hack and test on your feature (on a branch)
3. Document it
4. Submit a PR
5. Sign the CLA if necessary
6. Keep polishing it until received enough LGTM
7. Profit!

## Pull request requirements

For a pull request to be considered at all it has to meet these requirements:

1. Regardless if the code introduces new features or fixes bugs or regressions, it must have comprehensive tests.
1. The code must be well documented (see the ‘Documentation’ section below).
1. The commit messages must properly describe the changes, see further below.
1. All Asynkron projects must include Asynkron copyright notices. Each project can choose between one of two approaches:

   1. All source files in the project must have a Asynkron copyright notice in the file header.
   1. The Notices file for the project includes the Asynkron copyright notice and no other files contain copyright
      notices. See http://www.apache.org/legal/src-headers.html for instructions for managing this approach for
      copyrights.

   Proto.Actor uses the first choice, having copyright notices in every file header.

### Additional guidelines

Some additional guidelines regarding source code are:

- Files should start with a `Copyright (C) 2017 Asynkron.se <http://www.asynkron.se>` copyright header.
- Keep the code [DRY](http://programmer.97things.oreilly.com/wiki/index.php/Don%27t_Repeat_Yourself).
- Apply the [Boy Scout Rule](http://programmer.97things.oreilly.com/wiki/index.php/The_Boy_Scout_Rule) whenever you have
  the chance to.
- Never delete or change existing copyright notices, just add additional info.
- Do not use `@author` tags since it does not
  encourage [Collective Code Ownership](http://www.extremeprogramming.org/rules/collective.html).
  - Contributors , each project should make sure that the contributors gets the credit they deserve—in a text file or
    page on the project website and in the release notes etc.

If these requirements are not met then the code should **not** be merged into dev or master, or even reviewed -
regardless of how good or important it is. No exceptions.

Whether or not a pull request (or parts of it) shall be back- or forward-ported will be discussed on the pull request
discussion page, it shall therefore not be part of the commit messages. If desired the intent can be expressed in the
pull request description.

## Documentation

All documentation must abide by the following maxims:

- Example code should be run as part of an automated test suite.
- Generation should be **completely automated** and available for scripting.

All documentation must be in markup format.

### Go doc

Proto.Actor generates Go-style API documentation using `go doc`.

## External dependencies

All the external runtime dependencies for the project, including transitive dependencies, must have an open source
license that is equal to, or compatible with, [Apache 2](http://www.apache.org/licenses/LICENSE-2.0).

This must be ensured by manually verifying the license for all the dependencies for the project:

1. Whenever a committer to the project changes a version of a dependency in the build file.
2. Whenever a committer to the project adds a new dependency.
3. Whenever a new release is cut (public or private for a customer).

Which licenses are compatible with Apache 2 are defined
in [this doc](http://www.apache.org/legal/3party.html#category-a), where you can see that the licenses that are listed
under `Category A` are automatically compatible with Apache 2, while the ones listed under `Category B` need
additional action:

> Each license in this category requires some degree
> of [reciprocity](http://www.apache.org/legal/3party.html#define-reciprocal); therefore, additional action must be taken
> in order to minimize the chance that a user of an Apache product will create a derivative work of a
> reciprocally-licensed portion of an Apache product without being aware of the applicable requirements.

Each project must also create and maintain a list of all dependencies and their licenses, including all their transitive
dependencies. This can be done either in the documentation or in the build file next to each dependency.

## Creating commits and writing commit messages

Follow these guidelines when creating public commits and writing commit messages.

1. If your work spans multiple local commits (for example; if you do safe point commits while working in a feature
   branch or work in a branch for a long time doing merges/rebases etc.) then please do not commit it all but rewrite
   the history by squashing the commits into a single big commit which you write a good commit message for (like
   discussed in the following sections). For more info read this
   article: [Git Workflow](http://sandofsky.com/blog/git-workflow.html). Every commit should be able to be used in
   isolation, cherry picked etc.

2. The first line should be a descriptive sentence what the commit is doing, including the ticket number. It should be
   possible to fully understand what the commit does—but not necessarily how it does it—by just reading this single
   line. We follow the “imperative present tense” style for commit
   messages ([more info here](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)).

   It is **not ok** to only list the ticket number, type "minor fix" or similar.
   If the commit is a small fix, then you are done. If not, go to 3.

3. Following the single line description should be a blank line followed by an enumerated list with the details of the
   commit.

4. You can request review by a specific team member for your commit (depending on the degree of automation we reach, the
   list may change over time):
   - `Review by @gituser` - if you want to notify someone on the team. The others can, and are encouraged to
     participate.

Example:

    enable Travis CI #1

    * Details 1
    * Details 2
    * Details 3

# Supporting infrastructure

## Continuous integration

Each project should be configured to use a continuous integration (CI) tool (i.e. a build server à la Travis).

## Related links

- [Proto.Actor Contributor License Agreement](https://cla-assistant.io/keecon/protoactor-go)
- [Proto.Actor Issue Tracker](https://github.com/keecon/protoactor-go/issues)
