<script>
  import { daysAbbr, months } from '$lib/constants';
  import { addDays, getWeek, isDate, isSameDay, startOfWeek } from 'date-fns';
  import { onMount } from 'svelte';
  import Clock from './clock.svelte';
  import Icon from './icon.svelte';
  import Modal from './modal.svelte';

  export let value;
  export let show = false;

  const rows = [ ...Array(6).keys() ];
  const cols = [ ...Array(7).keys() ];
  const now = new Date();
  let year = now.getFullYear();
  let month = now.getMonth();
  let day = now.getDate();
  let hour = now.getHours();
  let minute = now.getMinutes();
  let second = now.getSeconds();
  let calendar = [];

  onMount(() => setDateToValues(now));
  $: setValueToDate(year, month, day, hour, minute, second);
  $: setDateToValues(value);

  function buildCalendar(y = year, m = month) {
    const date = new Date(y, m);
    let curDate = startOfWeek(date, { weekStartsOn: 1 });
    return rows.map(() => {
      const week = [];
      cols.forEach(() => {
        week.push(curDate);
        curDate = addDays(curDate, 1);
      });
      return week;
    });
  }

  function setValueToDate(date) {
    if (isDate(date)) {
      year = date.getFullYear();
      month = date.getMonth();
      day = date.getDate();
      hour = date.getHours();
      minute = date.getMinutes();
      value = new Date(year, month, day, hour, minute, second);
    }
    else {
      if (hour < 0) {
        day--; hour = 23;
      }
      else if (hour > 23) {
        day++; hour = 0;
      }
      else if (minute < 0) {
        hour--; minute = 59;
      }
      else if (minute > 59) {
        hour++; minute = 0;
      }
      else if (second < 0) {
        minute--; second = 59;
      }
      else if (second > 59) {
        minute++; second = 0;
      }
      value = new Date(year, month, day, hour, minute, second);
    }
    calendar = buildCalendar(year, month);
  }

  function setDateToValues(date) {
    year = date.getFullYear();
    month = date.getMonth();
    day = date.getDate();
    hour = date.getHours();
    minute = date.getMinutes();
    second = date.getSeconds();
    calendar = buildCalendar(year, month);
  }
</script>

<Modal width="700px" bind:show>
  <div class="datepicker">
    <div class="date">
      <div class="field">
        <input type="number" bind:value={day} />
        <select bind:value={month}>
          {#each months as monthName, index}
            <option value={index}>{monthName}</option>
          {/each}
        </select>
        <input type="number" bind:value={year} />
      </div>

      <table class="calendar">
        <thead>
          <tr>
            <th></th>
            {#each daysAbbr as dayName}
              <th>{dayName}</th>
            {/each}
          </tr>
        </thead>

        {#each calendar as week}
          <tr>
            <td class="week">{getWeek(week[0])}</td>
            {#each week as day}
              <td class="day">
                <button
                  on:click={() => setValueToDate(day)}
                  type="button"
                  class="button-small"
                  class:active={isSameDay(value, day)}
                  class:notinmonth={day.getMonth() !== month}
                >{day.getDate()}</button>
              </td>
            {/each}
          </tr>
        {/each}
      </table>
    </div>

    <div class="time">
      <div class="field">
        <input type="number" bind:value={hour} placeholder="hours" />
        <span class="label">:</span>
        <input type="number" bind:value={minute} placeholder="mins" />
        <span class="label">:</span>
        <input type="number" bind:value={second} placeholder="secs" />
      </div>
      <Clock date={value} />
    </div>
  </div>

  <div slot="footer" class="footer">
    <button class="btn secondary" type="button" on:click={() => value = new Date()}>
      <Icon name="o" /> Set to now
    </button>
    <button class="btn" type="button" on:click={() => show = false}>
      <Icon name="check" /> OK
    </button>
  </div>
</Modal>

<style>
  .datepicker {
    display: grid;
    grid-template: 1fr / 1fr 1fr;
    gap: 0.5rem;
  }

  .calendar {
    width: 100%;
    text-align: center;
  }
  .calendar thead th {
    opacity: 0.5;
    padding-top: 8px;
    padding-bottom: 8px;
  }
  .calendar .week {
    text-align: right;
    opacity: 0.5;
    padding-right: 10px;
  }
  .calendar .day button {
    display: block;
    width: 100%;
    padding-top: 8px;
    padding-bottom: 8px;
  }
  .calendar .day button.active {
    background-color: #00008b;
    color: #fff;
  }
  .calendar .day button.notinmonth {
    opacity: 0.6;
  }

  .time {
    display: flex;
    flex-direction: column;
  }
  .time input {
    text-align: center;
  }
  .time :global(.clock) {
    height: 150px;
    width: 150px;
    margin: auto;
  }

  .footer {
    display: flex;
    gap: 0.5rem;
    justify-content: right;
  }
</style>
